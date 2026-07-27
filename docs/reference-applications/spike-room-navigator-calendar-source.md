# Spike: non-Microsoft booking source on Webex Room Navigator

Status: Draft
Owner: Workplace Orchestrator Product

Spike run: 2026-07-25. All findings reflect documentation as published on that date.

## Question

Can a Cisco Webex Room Navigator (and the associated Webex room device / workspace) display and accept room reservations that originate from a source that is **not** Microsoft 365/Exchange and **not** Google Workspace — for example a self-hosted Nextcloud calendar, a CalDAV server, or a custom application acting as the booking system of record?

Sub-questions: (1) what mechanisms exist; (2) which are supported and can carry a non-MS/non-Google source of truth; (3) can the device originate an ad-hoc booking that flows back to the external system; (4) is occupancy/people-count available without Cisco Spaces; (5) what is the vendor-withdrawal risk.

## Verdict

**FEASIBLE WITH CONSTRAINTS — and the binding constraint is legal, not technical.**

Technically this is a solved problem and Cisco documents it: `xCommand Bookings Put` replaces the device's stored booking list with arbitrary JSON from any external source, and the `Bookings` xAPI additionally emits `BookingRequested` / `ExtensionRequested` / `SyncRequired` events so an external system can be the authority for on-device "Book now" too. Hybrid Calendar is irrelevant to this path (it only supports Exchange, Microsoft 365 and Google Calendar), and none of the booking xAPI is gated on it.

The constraint that matters: Cisco's RoomOS release notes attach special terms to exactly this API. Using `xCommand Bookings Put` (or any other API achieving the same result) to provide a "Scheduled Meeting Join Experience" is permitted **only for non-commercial use** — defined as internal business operations, explicitly *not* as part of an income-generating product — unless Cisco grants separate written permission. A commercial Workplace Orchestrator product sits squarely in the category that requires Cisco's permission, and Cisco reserves the right to revoke API access at its sole discretion. **The go/no-go here is a conversation with `devsupport@webex.com`, not an engineering test.**

## Mechanisms assessed

| Mechanism | Supported? | Can a non-MS/non-Google system be source of truth? | Key constraints | Source link |
|---|---|---|---|---|
| Webex Hybrid Calendar | Yes, fully supported and productised | **No.** Only Microsoft Exchange (on-prem, via Expressway connector), Microsoft 365, and Google Calendar. No CalDAV, no ICS feed, no generic connector documented. | Cannot deploy Google and Microsoft cloud back-ends in the same Webex org. Exchange path needs an Expressway-C cluster and an impersonation account. | [Deployment guide for Hybrid Calendar](https://help.webex.com/en-us/article/n6cwujdb/Deployment-guide-for-Hybrid-Calendar) (page dated 14 Jul 2026) |
| `xCommand Bookings Put` (device xAPI, local or via cloud xAPI) | Yes, publicly documented — **but carries special terms and conditions** | **Yes.** "Replaces the list of stored bookings… The meeting information is provided in JSON format." Source is entirely caller-defined. | "Applies to: All products". Requires user role ADMIN. Required JSON fields: `Id`, `Title`, `Number`, `Protocol`, `Organizer/Name`, `Time/StartTime`, `Time/Duration`. Commercial-use permission required (see Vendor-withdrawal risk). It is a *replace*, not a merge — caller owns full-list state. **Does not apply to a customer-managed stand-alone Room Navigator.** | API Reference Guide D15502.11, RoomOS 11.27, 05-2025, p.339; [roomos.cisco.com](https://roomos.cisco.com/xapi/Command.Bookings.Put/) |
| Cloud xAPI REST (`webexapis.com/v1/xapi/...`) | Yes, documented public API | Yes — it is the remote transport for `Bookings Put` and for reading statuses | Scopes `spark:xapi_commands` (commands) and `spark:xapi_statuses` (statuses). Needs a `deviceId` from the Devices API. Non-admin users and bots must be granted per-**workspace** API access in Control Hub (Workspaces → device card → Edit API Access; Full Access vs Read Only). Bots are single-org. | [xAPI reference](https://developer.webex.com/docs/api/v1/xapi); [Device Developers Guide](https://developer.webex.com/docs/api/guides/device-developers-guide) |
| Workspace Integrations | Yes, documented framework, public or private | Yes — and this is the only way to *receive* booking events from the cloud | Requires "Webex or Webex Edge registered endpoints". Runs as a machine account, managed in Control Hub → Workspaces → Integrations. Manifest declares fine-grained xAPI access rules per command/event/status (`xapiAccess.status` / `.commands` / `.events`, with wildcards permitted). Webhook **or** long polling, never both at once. Public integrations need Cisco approval; `org_private` ones are self-service. "Analytics, Meeting and XSI scopes do not currently work with Workspace Integrations." Notification subset is limited (see below). | [Workspace Integrations getting started](https://developer.webex.com/docs/workspace-integrations); [Technical Details](https://developer.webex.com/docs/workspace-integration-technical-details) |
| Macros (on-device JavaScript) | Yes, supported | Yes, in principle — a macro can call `Bookings Put` locally and use `HttpClient` to poll an external service | Same commercial-use terms apply. Per-device deployment and lifecycle. Cisco's own integration guide positions macros as the lowest-isolation option. *I did not verify macro `HttpClient` behaviour against a primary source in this spike — treat as unconfirmed.* | [RoomOS Integrations overview](https://roomos.cisco.com/doc/TechDocs/Integrations) (JS-rendered; content not directly extracted — see gaps) |
| Persistent Web App on Room Navigator | Yes, fully supported, and this is Cisco's **explicitly recommended** answer to "third party room booking" | **Yes.** Cisco's own use-case doc titled "Third party room booking" recommends exactly this. | Replaces the entire RoomOS UI on the Navigator and "can't be dismissed by end-users" — you lose the native booking and meeting-join UI and must rebuild it. Requires allowed-hosts and WebSocket/API-key config. Web app gets an embedded JSXAPI object with a **whitelisted** xAPI subset that includes `Bookings Book/Clear/Delete/Get/List/Put/Respond`, `Bookings Availability Status`, LED control, and `RoomAnalytics PeopleCount Current/Capacity` and `PeoplePresence`. | [Third party room booking use case](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/UseCases/CustomCalendaring.md); [WebAppsOnNavigator.md](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/Features/WebAppsOnNavigator.md); [Configure persistent web apps](https://help.webex.com/en-us/article/ohq3u6/Configure-persistent-web-apps-on-a-Room-Navigator) (page dated 10 Jul 2024) |
| On-device "local bookings" (no Hybrid Calendar) | Yes | Partly — the device itself becomes the store, which is *not* an external source of truth | Cisco documents that local bookings stored on the device are supported when Hybrid Calendar isn't enabled, and that ad-hoc booking is "provided by the on-device calendar or the selected calendar provider in Control Hub". Useful as a fallback, not as an integration. | [Enable room booking for shared mode devices](https://help.webex.com/en-us/article/nzp31jb/Enable-room-booking-for-shared-mode-Board-and-Room-Series-devices) (page dated 21 Apr 2026); [xConfiguration Bookings AdhocBooking Enabled](https://roomos.cisco.com/xapi/Configuration.Bookings.AdhocBooking.Enabled/) |
| On-prem / CUCM-registered devices | Partial | Yes via **local** xAPI; cloud xAPI only if linked to Webex Edge for Devices | `Bookings Put` is documented as "Applies to: All products" with Back-end "Any", so local xAPI works regardless of registration. Cloud xAPI and Workspace Integrations require Webex-registered or Webex Edge for Devices linked endpoints. Note `Bookings Respond` *is* explicitly limited to cloud-registered or Edge-linked devices — the restriction is per-command, not blanket. | [Webex Edge for Devices](https://help.webex.com/en-us/article/cy2l2z/Webex-Edge-for-Devices); [xCommand Bookings Respond](https://roomos.cisco.com/xapi/Command.Bookings.Respond/) |
| Stand-alone Room Navigator, **Webex-registered** (no room device) | Yes | **Yes.** The stand-alone Navigator has its own public API including `Bookings Book / Clear / Delete / Get / List / NotificationSnooze / Put / Respond` and `xStatus Bookings Availability Status` | A panel-only room display driven entirely by your backend. `Bookings Availability Status` returns `Free` / `FreeUntil` / `BookedUntil` with the boundary in `Bookings Availability TimeStamp`. No camera, so no people count. | Stand-alone Room Navigator API Reference Guide D15512.02, RoomOS 11.9, 04-2024, pp.31–34, 54 |
| Stand-alone Room Navigator, **customer managed** (`xConfiguration Provisioning Mode: Off`) | Yes, documented | **Only via a persistent web app** — the Bookings xAPI is unavailable | Fully cloud-independent: managed by SSH/local xAPI, customer owns upgrades and config, requires RoomOS 11.8 (Sept 2023) or later. **Every `Bookings` command and status in the stand-alone guide is annotated "Doesn't apply for a customer managed Room Navigator."** So this mode cannot use `Bookings Put` at all — you must render your own booking UI as a persistent web app. Correspondingly, the `Bookings Put` commercial-use terms do not bite here. | Stand-alone Room Navigator API Reference Guide D15512.02, RoomOS 11.9, 04-2024, pp.32–34, 54, 73–74 |

Products listed as supporting the `Bookings` commands include **Navigator**, Board Pro 55/75 (and G2), Desk Series, Codec EQ/Pro (and G2), Room 70 variants, Room Bar and Room Bar Pro. Note `Bookings Put` is marked **not** supported on devices in Microsoft Teams Rooms (MTR) mode.

## Ad-hoc booking

**Yes — a full round trip is documented, and this is the strongest finding of the spike.**

The `Bookings` xAPI is not write-only. It exposes a request/response pattern that lets an external system remain the authority even when the user books at the panel:

- `xEvent Bookings BookingRequested` fires with `BookingRequestUUID`, `StartTime`, `Duration`, `Title` when a user requests a booking on the device. This is the hook an external booking system needs.
- Companion events exist for the rest of the lifecycle: `ExtensionRequested`, `EditRequested`, `BookingCreated`, `BookingFailed`, `ExtensionFailed`, `EditFailed`, `Deleted`, `CheckedIn`, `BookingMoved`, `Updated`, `Start`, `End`, `TimeRemaining`.
- `xEvent Bookings SyncRequired` signals that the device wants its booking list refreshed — the natural trigger to re-push via `Bookings Put`.
- `xCommand Bookings Book` books the room from the device or via API, accepting `BookingRequestUUID`, `StartTime`, `Duration`, `Title`, `Organizer{Name,Email,Id}`, `Number`, `Protocol`, and `MeetingPlatform` (MicrosoftTeams / GoogleMeet / Zoom / Webex / Other / None). Default duration 30 minutes. It "is only available for devices that support the room scheduling feature, refer to the AdhocBooking Enabled setting."
- `xCommand Bookings Extend`, `Delete`, `Edit`, `CheckIn`, `CheckOut`, `NotificationSnooze` cover the rest of the panel interactions. Check-in behaviour is configurable via `Bookings CheckIn Enabled`, `WindowDuration` and `AutoCheckInDuration`.
- `xCommand Bookings Respond` (Accept/Decline with optional comment) exists but is for **responding to a meeting invitation** addressed to the room, not for acknowledging an ad-hoc request. It is documented as available only for cloud-registered devices or on-prem devices linked to Webex Edge for Devices.

Important caveat on the cloud path: both `Bookings BookingRequested` and `Bookings SyncRequired` are annotated on roomos.cisco.com as **"Only available with Workspace Integrations"** for cloud delivery. So the round trip works remotely only if you build a Workspace Integration — plain cloud xAPI polling will not deliver these events.

Ad-hoc booking must be explicitly enabled: `xConfiguration Bookings AdhocBooking Enabled` defaults to **False**. Related knobs: `Bookings AdhocBooking BookAhead Enabled`, `Bookings AdhocBooking BookingTimeout`, `Bookings MaxDuration`, `Bookings AllowDecline`, `Bookings DeclinableBookingType`.

Cisco's docs state that ad-hoc booking "is provided by the on-device calendar or the selected calendar provider in Control Hub" — meaning the *native* UI's ad-hoc booking writes either to the device or to the Control Hub-configured (i.e. Microsoft/Google) provider. Routing it to a third system is what `BookingRequested` is for. **Inference (not directly stated in any source I found): the intended pattern is (a) intercept `BookingRequested`, (b) create the reservation in your own system, (c) push the refreshed list with `Bookings Put`.** I could not find a Cisco document that spells this sequence out end to end, and I could not confirm whether the device shows a pending/spinner state while waiting, or how long it waits before timing out (`Bookings AdhocBooking BookingTimeout` is presumably the control, but I did not read its description).

For the native wall-panel flow with a supported calendar back-end, Cisco documents that a booking made at the panel "becomes visible in the meeting room's calendar, but it doesn't show up in your personal calendar", bookings are limited to the present time until midnight the same day, and minimum duration is 10 minutes.

## Occupancy data without Cisco Spaces

**Yes — people count and presence are first-class device xAPI statuses, readable locally and via cloud xAPI, with no Cisco Spaces dependency anywhere in the documentation.**

Cisco's Workspace Integrations overview lists, as example use cases in its own words, both an office map showing rooms "occupied or vacant, based on actual people presence from the video device's sensors and AI" and — verbatim — **"Un-book a booked meeting room when it is detected that nobody showed up."** That is precisely the no-show release feature, named by Cisco as an intended use of the free, documented integration framework.

Cisco Spaces is a **downstream consumer** of this telemetry, not its source and not a gate on it. Cisco Spaces' own runbook describes Webex devices as sending telemetry "on changes to sensor values" into Spaces; what Spaces adds is cross-portfolio aggregation, floor heatmaps, occupancy playback and fusion with Wi-Fi/BLE sensors. A Webex blog post (published 10 Jun 2020, updated 21 Dec 2023) states that Webex Room, Board and Desk devices "come with intelligent people counting sensors directly embedded without any additional cost", and help.webex.com's native workspace-metrics articles do not mention Cisco Spaces at all. **Note this is a well-supported inference, not a Cisco licensing statement — I found no page that formally sets out "included vs. requires Spaces".**

Verified statuses and configs:

- `xStatus RoomAnalytics PeopleCount Current` — head-detection count. Cloud-readable: `GET https://webexapis.com/v1/xapi/status/?deviceId=...&name=RoomAnalytics.PeopleCount.Current`. Roles Admin/Integrator/User. Marked privacy-impacting.
- `xStatus RoomAnalytics PeopleCount Capacity` — configurable in Control Hub.
- `xStatus RoomAnalytics PeoplePresence`, `xStatus RoomAnalytics UltrasoundPresence`, `xStatus RoomAnalytics RoomInUse`, `xStatus RoomAnalytics OccupancySensor[n] RoomInUse`, `xStatus RoomAnalytics Engagement CloseProximity`.
- `xConfiguration RoomAnalytics PeopleCountOutOfCall` — must be **On** to count people outside a call; default behaviour ("Off") counts only in-call or with self-view on.
- `xConfiguration RoomAnalytics PeoplePresenceDetector` — enables the ultrasound-based detector.
- `xConfiguration RoomAnalytics PeoplePresence Input HeadDetector` / `Input Ultrasound`.

Additional verified nodes: `xStatus RoomAnalytics RoomInUse` (True/False, "the final effective room-in-use state" derived from presence, call state and sharing state — and notably **not** flagged privacy-impacting, unlike `PeopleCount Current`), and `xCommand RoomAnalytics RoomInUse Set`, which lets an external sensor feed occupancy in as an additional input. `RoomInUse` is likely the better signal for a no-show feature than raw people count: it is coarser, privacy-clean, and already fuses the inputs.

`xStatus RoomAnalytics PeoplePresence` returns Yes/No/**Unknown** (Unknown when `PeoplePresenceDetector` is Off). Cisco notes presence updates immediately on entry but **"after the room becomes empty, it may take up to two minutes for the status to change"** — a floor on how fast any release decision can be made.

### RoomOS may already ship this feature

Before building a no-show release loop, note that RoomOS has one natively, with no macro, no integration and no Cisco Spaces:

- `xConfiguration Bookings CheckIn Enabled` (default **False**) — Cisco's own rationale: "In a significant percentage of meetings no one shows up and the meeting resources are wasted during the booking… If no one arrives for the booking, the meeting room can be released." Documented as available on shared bookable devices **and standalone Room Navigator**.
- `xConfiguration Bookings CheckIn WindowDuration` (5/10/15/20/30/60 min, default 10) — "If no one checks in during this period, manually or automatically, the room is automatically released."
- `xConfiguration Bookings CheckIn AutoCheckInDuration` (8–24h, default 8) — threshold for auto-check-in on long/all-day bookings.

Automatic check-in triggers on ad-hoc booking, call join, share session, or when people count reaches 1+, and starts four minutes after the booking begins. The documented config combination is `RoomAnalytics PeopleCountOutOfCall: On` **and** `Standby WakeupOnMotionDetection: On`.

**This should be evaluated before building anything custom.** The open question is what "released" means when the booking's source of truth is *your* system rather than Exchange/Google — the device presumably emits a `Bookings Deleted` / `Updated` event and expects the backend to follow, but I did not verify that. (**Inference, unverified.**)

Three caveats that directly threaten a naive custom no-show implementation:

1. **Standby returns -1.** Cisco states that when the device is in standby, or the feature is off, `PeopleCount Current` returns `-1`. An empty room will often be in standby, so `-1` must be treated as "unknown", and distinguishing "empty" from "asleep" needs `PeoplePresence` / `RoomInUse` or `xConfiguration Standby WakeupOnMotionDetection`.
2. **People count is never active in networked standby** (the deeper low-power stage on Desk, Desk Mini, Room Bar, Room Bar Pro). Cisco notes those devices will not enter networked standby until people count reaches 0 — useful, but it means the count is frozen once asleep.
3. **Accuracy.** Cisco warns objects or pictures may be counted as heads and that a minimum two-minute duration is needed for a reliable average.

Hardware dependency: for Codec Plus the status requires a Quad Camera; for Codec Pro a Quad Camera or SpeakerTrack 60; for Codec EQ a Quad Camera. Desk, Desk Mini, Desk Pro and Desk Pro G2 support out-of-call counting natively. People **presence** can alternatively use ultrasound (`PeoplePresence Input Ultrasound`), which needs no camera and is less privacy-sensitive but coarser and subject to the ~2-minute clear-out lag.

The Room Navigator itself is a touch panel with **no camera** and generates no people count of its own. Cisco's own wording is that the Navigator shows people count and presence "if that is available from the device it's connected to" — the Navigator appears in the `PeopleCount Current` product list because it relays and displays, not because it senses. (It does carry genuine environmental sensors, exposed as `Peripherals ConnectedDevice[n] RoomAnalytics AirQuality Index`, `AmbientTemperature`, `RelativeHumidity`.) **A Navigator-only room therefore has no native occupancy signal** — the no-show feature requires a paired room device with a supported camera, or an external sensor fed in via `xCommand RoomAnalytics RoomInUse Set`.

Supporting evidence that this is the blessed path: Cisco's **own example Workspace Integration manifest** requests `RoomAnalytics.*` and `Peripherals.ConnectedDevice[*].RoomAnalytics.*` as required status access, alongside `Standby.State` and `SystemUnit.State.NumberOfActiveCalls` — i.e. exactly the signal set a no-show detector needs, including the standby state required to disambiguate the `-1` case, and including the Room Navigator's own peripheral sensors. No Cisco Spaces scope appears anywhere in the manifest schema.

Delivery constraints for Workspace Integrations — three that matter:

1. Notifications cover "only a small subset of the statuses and events… such as room analytics, custom UI events, boot events and a few more", though all public statuses remain *queryable*. The published subscribable-status table puts `RoomAnalytics.PeoplePresence` and `RoomAnalytics.PeopleCount.Current` in the **"sent immediately when changed"** tier; `Bookings.Availability.Status` is throttled to at most once per minute and `RoomAnalytics.AmbientNoise.Level.A` to once per five minutes.
2. "Change notifications are only supported for shared devices (devices in workspaces). Personal devices are currently not supported."
3. **A tenant opt-in gates the data**: "the RoomAnalytics data is only sent if the customer has opted into Workspace Utilization data and/or Workspace Environmental data in Control Hub". This is a per-customer admin setting under Workspaces → Settings and is a real deployment prerequisite to document for customers. Whether the same opt-in also gates plain REST `GET /v1/xapi/status` reads of RoomAnalytics is **not documented and unverified**.

Delivery is webhook **or** HTTP long-polling a `pollUrl` (long poll returns at most 10 messages, waits up to 20s, 10-minute retention) — not both on one integration.

There is also a separate native **Workspace Metrics API** (`GET /v1/workspaceMetrics`, scope `spark-admin:workspace_metrics_read`) whose `metricName` enum includes `peopleCount`. It is historical/aggregated — `none`/`hourly` capped at a 48-hour span, `daily` at 30 days — and Cisco notes the mean "should be interpreted as the average number of people in the room when it is in use". Useful for tuning no-show thresholds, **not** for real-time release decisions.

Cloud query limits: at most **10 status expressions per request**; each device object carries a `permissions` array and "the bot or user can make xAPI calls only to the devices where 'xapi' is part of this array". No numeric polling rate limit is published — only a generic 429 with `Retry-After` — which is a real unknown if you intend to poll many rooms rather than subscribe.

## Vendor-withdrawal risk

**Moderate on the technical surface. High on the licensing surface. The risk is not that Cisco removes the API — it is that Cisco has pre-emptively reserved the right to forbid your commercial use of it.**

The RoomOS 26 release notes carry a section headed "Permitted Commercial Use for Scheduled Meeting Join Experience" stating that providing "a scheduled meeting join experience such as one button to push" requires separate permission for commercial use, and that this "includes use of any API that updates the device with calendar data from an external source to provide this functionality, including `xCommand Bookings Put` or previous private APIs such as `bookingsputxml`. In addition, using other APIs to accomplish the same functionality also requires permission for commercial use."

Non-commercial use is defined as "use solely for your internal business operations and not for activities that involve using the API as part of, or in furtherance of, an income-generating service or product, whether directly or indirectly." Cisco adds that it "reserves the right to revoke your license to use our API if, in our sole discretion, we deem that your use is for unauthorized commercial purposes", and directs questions to `devsupport@webex.com`.

The same terms appear, near-verbatim, in the stand-alone Room Navigator API Reference Guide (D15512.02, RoomOS 11.9, dated 04-2024) under "Permitted commercial use for scheduled meeting join experience", and the command entry itself in the main API Reference Guide (D15502.11, RoomOS 11.27, 05-2025, p.339) carries: "This API has special terms and conditions, please refer to the 'About the API – Terms and Conditions' section in the API guide."

That the identical clause appears in an April 2024 guide, a May 2025 guide, and the current RoomOS 26 release notes means this is **settled, deliberate Cisco policy, not a transitional notice**. It is not going to be quietly dropped.

Reading of this, marked as **inference**: the restriction targets the *scheduled meeting join experience* (one-button-to-push), not room availability display per se. A product that pushes bookings purely to drive the Navigator's availability panel and booking UI, without offering a one-touch join of a scheduled meeting, is arguably outside the restricted use case. That reading is **not confirmed** and is exactly the kind of thing that needs Cisco in writing — the sentence "using other APIs to accomplish the same functionality also requires permission" is drafted broadly enough to be read either way.

Mitigating factors on the technical side:
- The `Bookings` xAPI surface is large, coherent, versioned, and present in the current RoomOS 26 generation — not a vestigial corner. It has *grown* (check-in/check-out, book-ahead, decline policy) rather than shrunk.
- Cisco explicitly names the anti-no-show and occupancy-map use cases in its own developer marketing.
- Cisco documents the "Third party room booking" persistent-web-app path as a supported product capability, and the Navigator data sheet positions third-party apps as protecting existing room-booking investments.
- The escape hatch is genuine and now confirmed: a **customer-managed** stand-alone Room Navigator has no access to the `Bookings` xAPI at all, yet fully supports a persistent web app rendering your own booking UI against your own backend. That configuration is outside the scope of the `Bookings Put` terms entirely, needs no Webex cloud, and needs no Cisco permission — at the cost of the native RoomOS UI and of owning software upgrades yourself.

Aggravating factors:
- Cisco has been consolidating workspace/occupancy analytics into Cisco Spaces and Control Hub. The Cisco Spaces ↔ Control Hub integration is documented, but a targeted search across 2024–2026 release notes and "What's new in Control Hub" found **no deprecation notice for any native RoomOS occupancy API**, and all `RoomAnalytics` nodes remain live in the current xAPI reference. One forward-looking signal is worth tracking, though it cuts the other way: the Cisco Spaces device compatibility matrix notes that "Webex Room Occupancy will be changing from xStatus RoomAnalytics PeopleCount Current to xStatus RoomAnalytics RoomInUse to cover more scenarios" — that is a change to which node *Spaces consumes*, not a deprecation of `PeopleCount Current`, and it is a hint that `RoomInUse` is the more durable node to build on.
- `Bookings Put` explicitly supersedes "previous private APIs such as `bookingsputxml`", which is direct evidence Cisco has already once migrated this capability from a private API. That cuts both ways: it shows willingness to churn, and willingness to make it public.

## What I could not determine

- **Whether Cisco would grant commercial permission** for a product like Workplace Orchestrator, on what terms, at what price, or how long the process takes. This is the single largest unknown and it is not answerable from public sources.
- **Whether the "Scheduled Meeting Join Experience" restriction bites** if the product never offers one-button-to-push. Unresolved; the drafting is ambiguous.
- The release notes say to "comply with the permitted commercial use terms below", but **no such terms appear below that paragraph** in the RoomOS 26 release notes I retrieved. Either the terms live elsewhere or the document is incomplete. I could not locate them.
- **Whether a Workspace Integration manifest may declare `Bookings Put`** in its xAPI access rules. The manifest schema takes a free-form `xapiAccess.commands` array of `path` entries (Cisco's own example uses `Message.Send`, and wildcards such as `RoomAnalytics.*` are shown for statuses), so `Bookings.Put` should be declarable — but I found no explicit confirmation, and no published exclusion list. **Inference, unverified.**
- **Whether the cloud xAPI restricts multi-line commands.** `Bookings Put` is a multi-line command with a raw payload; roomos.cisco.com shows `POST https://webexapis.com/v1/xapi/command/Bookings.Put` with a `body` field, so it appears supported, but I did not find prose confirming multi-line/large-payload handling or size limits.
- **Persistence and lifetime of bookings pushed via `Bookings Put`**: does the list survive reboot, does it expire, how often must it be re-pushed, does anything else (Hybrid Calendar sync) overwrite it? Not documented in anything I read. `Bookings SyncRequired` implies the device asks when it needs a refresh, but the contract is unstated.
- **Whether `Bookings Put` and a configured Hybrid Calendar can coexist**, and which wins. Untested and undocumented as far as I found.
- **Cloud xAPI rate limits** — no published numeric limit; only a generic 429 with `Retry-After`. Material if the design polls rather than subscribes.
- **Whether the Control Hub "Workspace Utilization data" opt-in gates REST status reads** of RoomAnalytics, or only Workspace Integration notifications. Documented only for the latter.
- **What "release the room" does when your system owns the booking.** The native `Bookings CheckIn` auto-release presumably emits an event your backend must honour, but I did not verify which event, or whether auto-release works at all when bookings arrived via `Bookings Put` rather than from a Control Hub calendar provider. This is arguably the single most important *technical* unknown for the product.
- **Minimum RoomOS versions per xAPI node.** roomos.cisco.com does not display an "introduced in" field, and I read only the RoomOS 11.9, 11.27 and 26 guides. Version floors for `Bookings CheckIn`, `BookingRequested` and `AdhocBooking BookAhead` are unconfirmed.
- **Whether a paid Cisco Spaces tier is needed for portfolio-level dashboards** of Webex-device-only occupancy (as opposed to per-device API access, which is clearly free).
- **Macro `HttpClient` specifics** (outbound calls to a CalDAV/REST server from the device, TLS/cert requirements, `HttpClient Mode` and `AllowInsecureHTTPS`). Not verified; roomos.cisco.com's Integrations tech doc is JS-rendered and I did not extract it.
- **Licensing for room booking itself.** A search result asserted the native Cisco room-booking solution is licence-free when paired with a Cisco video device, but I could not confirm this against a primary Cisco source (the Navigator data sheets returned HTTP 403). Treat as unverified. Whether a "Cisco Room Scheduler" SKU or Webex Device licence tier is a prerequisite for the `Bookings` xAPI is likewise unconfirmed.
- **Whether a stand-alone Webex-registered Navigator can be driven by `Bookings Put` without any Hybrid Calendar configured at all**, or whether Control Hub still insists on a calendar provider being selected. Unverified.
- **Whether people count is obtainable in a Navigator-only room.** A stand-alone Navigator has no camera, so head-detection people count is presumably unavailable there (**inference**) — the no-show feature likely requires a paired room device with a supported camera. Not explicitly stated in any source I read.
- **On-prem/TMS/TMSXE booking on a Navigator** — the relevant Cisco Community thread returned HTTP 403. Unverified.
- **`Bookings BookingTimeout` semantics**, and the device's UI behaviour while an external system decides on a booking request.
- Note on method: community.cisco.com and the Navigator data sheets returned HTTP 403 to automated retrieval and could not be read. The two API Reference Guide PDFs *were* eventually obtained in full and read directly, so the commercial-use terms and the `Bookings Put` specification in this document are quoted from primary Cisco PDFs, not from search snippets. No retrieved page contained text addressed to an AI agent.
- Minor signal on documentation maturity, noted for what it is worth: the stand-alone Navigator guide (D15512.02) contains an unfilled editorial placeholder — "You can download the software for a customer managed Room Navigator from: WHERE?" — in the customer-managed setup section. The customer-managed path is real but evidently lightly trafficked.

## Recommended next step

**Two things, in parallel, in this order of importance.**

1. **Email `devsupport@webex.com` describing the intended commercial use** — a product that pushes reservations from a customer's own booking system (Nextcloud/CalDAV/custom) onto Webex room devices via `xCommand Bookings Put`, and asking (a) whether this falls under the "Scheduled Meeting Join Experience" restriction if no one-button-to-push join is offered, (b) what the permission process and terms are, and (c) where the "permitted commercial use terms" referenced in the RoomOS 26 release notes actually live. This costs one email and is the actual gate on the product. Do it before writing code.

2. **Cheapest technical proof, ~half a day, no hardware purchase:** request a Webex developer sandbox org (or use one Room Navigator paired to any Room Bar / Board on an existing tenant), then from a personal access token with `spark:xapi_commands` + `spark:xapi_statuses`:
   - `GET /v1/devices` to obtain a `deviceId`;
   - grant the account Full Access on that workspace in Control Hub (Workspaces → device card → Edit API Access);
   - `PATCH /v1/deviceConfigurations` to set `Bookings.AdhocBooking.Enabled = True` and `RoomAnalytics.PeopleCountOutOfCall = On`;
   - `POST /v1/xapi/command/Bookings.Put` with a two-entry JSON payload using only the required fields (`Id`, `Title`, `Number`, `Protocol`, `Organizer/Name`, `Time/StartTime`, `Time/Duration`);
   - photograph the Navigator and confirm the bookings render on the panel and home screen, then reboot the device and re-query to settle the persistence question;
   - `GET /v1/xapi/status/?name=RoomAnalytics.PeopleCount.Current` and `...=RoomAnalytics.RoomInUse` with the room empty and occupied, to confirm the `-1`/standby behaviour;
   - set `Bookings.CheckIn.Enabled = True` and let a `Bookings Put`-sourced booking go unattended past the check-in window, to see whether the **native** auto-release fires and what event it emits.

   Then, only if step 2 passes, build a throwaway **`org_private`** Workspace Integration subscribing to `Bookings BookingRequested` and `RoomAnalytics PeoplePresence` to prove the ad-hoc round trip and the no-show signal. Private integrations need no Cisco approval, so this is self-service — but first enable the Workspace Utilization data opt-in in Control Hub under Workspaces → Settings, or RoomAnalytics notifications will silently not arrive.

   This settles persistence, Hybrid-Calendar coexistence, multi-line payload handling over cloud xAPI, native-vs-custom no-show release, and the occupancy standby behaviour — the biggest technical unknowns — for the price of a sandbox request.

**Design note worth carrying into the next step:** prefer `RoomAnalytics RoomInUse` over `PeopleCount Current` as the occupancy signal (privacy-clean, already fuses inputs, and is the node Cisco Spaces is itself migrating toward), and prefer Workspace Integration subscriptions over polling.

## Sources

Every URL below was actually fetched during this spike unless marked otherwise.

- [xCommand Bookings Put — roomos.cisco.com](https://roomos.cisco.com/xapi/Command.Bookings.Put/) — no date shown (rendered via browser; JS app)
- [xCommand Bookings Book — roomos.cisco.com](https://roomos.cisco.com/xapi/Command.Bookings.Book/) — no date shown
- [xCommand Bookings Respond — roomos.cisco.com](https://roomos.cisco.com/xapi/Command.Bookings.Respond/) — no date shown
- [xEvent Bookings BookingRequested — roomos.cisco.com](https://roomos.cisco.com/xapi/Event.Bookings.BookingRequested/) — no date shown
- [xEvent Bookings SyncRequired — roomos.cisco.com](https://roomos.cisco.com/xapi/Event.Bookings.SyncRequired/) — no date shown
- [xConfiguration Bookings AdhocBooking Enabled — roomos.cisco.com](https://roomos.cisco.com/xapi/Configuration.Bookings.AdhocBooking.Enabled/) — no date shown
- [xStatus RoomAnalytics PeopleCount Current — roomos.cisco.com](https://roomos.cisco.com/xapi/Status.RoomAnalytics.PeopleCount.Current/) — no date shown
- [xConfiguration RoomAnalytics PeopleCountOutOfCall — roomos.cisco.com](https://roomos.cisco.com/xapi/Configuration.RoomAnalytics.PeopleCountOutOfCall/) — no date shown
- [RoomOS 26 release notes (source markdown, cisco-ce/roomos.cisco.com)](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/WhatsNew/ReleaseNotesRoomOS_26.md) — RoomOS 26; source of the verbatim commercial-use terms
- [Third party room booking use case (CustomCalendaring.md)](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/UseCases/CustomCalendaring.md) — no date shown
- [Space management use case (SpaceManagement.md)](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/UseCases/SpaceManagement.md) — no date shown
- [Web apps on Navigator (WebAppsOnNavigator.md)](https://raw.githubusercontent.com/cisco-ce/roomos.cisco.com/master/doc/Features/WebAppsOnNavigator.md) — no date shown; source of the persistent-web-app xAPI whitelist
- [Workspace Integrations — Getting Started, developer.webex.com](https://developer.webex.com/docs/workspace-integrations) — no date shown
- [Workspace Integrations — Technical Details, developer.webex.com](https://developer.webex.com/docs/workspace-integration-technical-details) — no date shown; source of the manifest schema and example
- [xAPI reference, developer.webex.com](https://developer.webex.com/docs/api/v1/xapi) — no date shown
- [Devices / Device Developers Guide, developer.webex.com](https://developer.webex.com/docs/api/guides/device-developers-guide) — no date shown
- [Deployment guide for Hybrid Calendar, help.webex.com](https://help.webex.com/en-us/article/n6cwujdb/Deployment-guide-for-Hybrid-Calendar) — page dated 14 Jul 2026
- [Enable room booking for shared mode Board and Room Series devices, help.webex.com](https://help.webex.com/en-us/article/nzp31jb/Enable-room-booking-for-shared-mode-Board-and-Room-Series-devices) — page dated 21 Apr 2026
- [Room Navigator as a standalone device, help.webex.com](https://help.webex.com/en-us/article/iq6aw6/Room-Navigator-as-a-standalone-device) — page dated 18 Aug 2025
- [Set up Room Navigator as a room booking device, help.webex.com](https://help.webex.com/en-us/article/55ypt4/Set-up-Room-Navigator-as-a-room-booking-device) — page dated 16 Oct 2025
- [Configure persistent web apps on a Room Navigator, help.webex.com](https://help.webex.com/en-us/article/ohq3u6/Configure-persistent-web-apps-on-a-Room-Navigator) — page dated 10 Jul 2024
- [Outside room booking with a wall mounted Room Navigator, help.webex.com](https://help.webex.com/en-us/article/me7z9t/Outside-room-booking-with-a-wall-mounted-Room-Navigator) — page dated 16 Dec 2024
- [Webex Edge for Devices, help.webex.com](https://help.webex.com/en-us/article/cy2l2z/Webex-Edge-for-Devices) — referenced via search results; not directly fetched
- [RoomOS Integrations overview, roomos.cisco.com](https://roomos.cisco.com/doc/TechDocs/Integrations) — fetch returned title only (JS-rendered); content not extracted
- **Cisco Room Navigator (stand-alone) RoomOS 11 API Reference Guide, D15512.02, RoomOS 11.9, 04-2024** — `https://www.cisco.com/c/dam/en/us/td/docs/telepresence/endpoint/room_navigator-stand-alone/api/stand-alone-room-navigator-api-guide-roomos119.pdf` — obtained and read in full (75 pp.); source of the verbatim commercial-use terms, the stand-alone Bookings command set, and the customer-managed mode section
- **Cisco collaboration devices RoomOS 11 API Reference Guide, D15502.11, RoomOS 11.27, 05-2025** — `https://www.cisco.com/c/dam/en/us/td/docs/telepresence/endpoint/roomos-111/api-reference-guide-roomos-111.pdf` — obtained and read (727 pp.); source of the `Bookings Put` "Applies to: All products" entry, p.339
- [Utilization and environmental metrics for workspaces, help.webex.com](https://help.webex.com/en-us/article/nc6od6r/Utilization-and-environmental-metrics-for-workspaces) — page dated 1 Jun 2026
- [Workspace insights in Control Hub, help.webex.com](https://help.webex.com/en-us/article/yuzsuw/Workspace-insights-in-Control-Hub) — page dated 5 Mar 2026
- [Workspace Metrics API, developer.webex.com](https://developer.webex.com/admin/docs/api/v1/workspace-metrics/workspace-metrics) — no date shown
- [xAPI Query Status, developer.webex.com](https://developer.webex.com/docs/api/v1/xapi/query-status) and [Execute Command](https://developer.webex.com/docs/api/v1/xapi/execute-command) — no dates shown
- [Cisco Spaces occupancy runbook](https://runbooks.ciscospaces.io/docs/cisco-spaces-occupancy-runbook-cisco-validated) and [device compatibility matrices](https://runbooks.ciscospaces.io/docs/cisco-spaces-device-compatibility-matrices) — no dates shown
- [Returning to work with intelligent room capacity, Webex blog](https://blog.webex.com/collaboration-devices/returning-to-work-with-intelligent-room-capacity/) — published 10 Jun 2020, updated 21 Dec 2023
- roomos.cisco.com xAPI pages additionally consulted: `Status.RoomAnalytics.PeopleCount.Capacity`, `Status.RoomAnalytics.PeoplePresence`, `Status.RoomAnalytics.RoomInUse`, `Configuration.RoomAnalytics.PeoplePresenceDetector`, `Configuration.Bookings.CheckIn.Enabled`, `Configuration.Bookings.CheckIn.WindowDuration`, `Configuration.Bookings.CheckIn.AutoCheckInDuration` — no dates shown

Retrieval failures (HTTP 403, content not obtained): the Room Navigator for Wall and for Table data sheets, and the Cisco Community thread on Navigator room booking with on-prem devices and TMS/TMSXE. Both API Reference Guide PDFs initially returned 403 to automated fetch but were subsequently obtained in full and read directly.
