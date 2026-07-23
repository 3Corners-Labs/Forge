package evidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinIOArtifactStore struct {
	client   *minio.Client
	bucket   string
	metadata *PostgresStore
	now      func() time.Time
	maxBytes int64
}

func NewMinIOArtifactStore(client *minio.Client, bucket string, metadata *PostgresStore, maxBytes int64, now func() time.Time) (*MinIOArtifactStore, error) {
	if client == nil || bucket == "" || metadata == nil || maxBytes < 1 || now == nil {
		return nil, errors.New("MinIO artifact configuration is invalid")
	}
	return &MinIOArtifactStore{client: client, bucket: bucket, metadata: metadata, maxBytes: maxBytes, now: now}, nil
}

type PutArtifactRequest struct {
	TenantID, ProjectID, DeliveryRunID string
	MediaType, Classification          string
	RetentionClass, CreatorCapability  string
	Bytes                              []byte
}

func (store *MinIOArtifactStore) Put(ctx context.Context, request PutArtifactRequest) (ArtifactMetadata, error) {
	if len(request.Bytes) == 0 || int64(len(request.Bytes)) > store.maxBytes ||
		!safeObjectSegment(request.TenantID) || !safeObjectSegment(request.ProjectID) ||
		!safeObjectSegment(request.DeliveryRunID) ||
		request.MediaType == "" || request.Classification == "" || request.RetentionClass == "" ||
		request.CreatorCapability == "" {
		return ArtifactMetadata{}, errors.New("artifact request is invalid")
	}
	sum := sha256.Sum256(request.Bytes)
	digest := "sha256:" + hex.EncodeToString(sum[:])
	prefix := request.TenantID + "/" + request.ProjectID + "/" + request.DeliveryRunID
	temporary := prefix + "/.tmp/" + strings.TrimPrefix(digest, "sha256:")
	final := prefix + "/sha256/" + strings.TrimPrefix(digest, "sha256:")
	if _, err := store.client.PutObject(ctx, store.bucket, temporary, bytes.NewReader(request.Bytes), int64(len(request.Bytes)), minio.PutObjectOptions{ContentType: request.MediaType}); err != nil {
		return ArtifactMetadata{}, err
	}
	defer store.client.RemoveObject(context.Background(), store.bucket, temporary, minio.RemoveObjectOptions{})
	object, err := store.client.GetObject(ctx, store.bucket, temporary, minio.GetObjectOptions{})
	if err != nil {
		return ArtifactMetadata{}, err
	}
	verified, err := io.ReadAll(io.LimitReader(object, store.maxBytes+1))
	object.Close()
	if err != nil || len(verified) != len(request.Bytes) || sha256.Sum256(verified) != sum {
		return ArtifactMetadata{}, errors.New("temporary artifact digest verification failed")
	}
	info, statErr := store.client.StatObject(ctx, store.bucket, final, minio.StatObjectOptions{})
	if statErr == nil {
		if info.Size != int64(len(request.Bytes)) {
			return ArtifactMetadata{}, errors.New("immutable artifact key already has different content")
		}
	} else {
		code := minio.ToErrorResponse(statErr).Code
		if code != "NoSuchKey" && code != "NoSuchObject" {
			return ArtifactMetadata{}, statErr
		}
		if _, err := store.client.CopyObject(ctx, minio.CopyDestOptions{Bucket: store.bucket, Object: final}, minio.CopySrcOptions{Bucket: store.bucket, Object: temporary}); err != nil {
			return ArtifactMetadata{}, err
		}
	}
	at := store.now().UTC()
	metadata := ArtifactMetadata{
		ObjectKey: final, TenantID: request.TenantID, ProjectID: request.ProjectID,
		DeliveryRunID: request.DeliveryRunID, Digest: digest, SizeBytes: int64(len(request.Bytes)),
		MediaType: request.MediaType, Classification: request.Classification,
		RetentionClass: request.RetentionClass, CreatorCapability: request.CreatorCapability,
		CreatedAt: at.Format(time.RFC3339Nano),
	}
	if err := store.metadata.CommitArtifact(ctx, metadata); err != nil {
		return ArtifactMetadata{}, err
	}
	return metadata, nil
}

func safeObjectSegment(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	for _, character := range value {
		if (character < 'a' || character > 'z') && (character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') && character != '-' && character != '_' {
			return false
		}
	}
	return true
}
