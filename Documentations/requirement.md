# Yojana Portal: Production-Grade Production Requirements & Architecture Specification

## Executive Overview
This document specifies the requirements, architectural improvements, and future development guidelines for the Yojana Portal (majhiGov) as it prepares for high-concurrency production serving **100,000+ active citizens**. 

The system has been analyzed from three distinct critical perspectives:
1. **The Citizen (End User)**: Ensuring maximum accessibility, simple workflows, seamless multi-lingual translation, and real-time visibility into application lifecycles.
2. **The Administrator (Ops & Government Officers)**: Streamlining scheme curation, eligibility rule adjustments, citizen cohort notifications, and rapid citizen verification/decision processes.
3. **The CTO (Technical Scalability & Security)**: Maintaining low latency, high throughput, robust security (Aadhaar & PII encryption), strict rate-limiting, and resilient auto-translation fallback systems...

---

## 👥 Persona 1: The Citizen (User Requirements)
Citizens in rural and semi-urban settings require ultra-low cognitive load interfaces that function smoothly under moderate cellular network conditions.

### 1. Dynamic Eligibility Checker (Smart Explorer)
* **Status Check**: Provide instant feedback on why a citizen is not eligible for a specific scheme, rather than a generic "Not Eligible" message.
* **Auto-Filled Cohorts**: Once registered and logged in, the portal must automatically suggest matches using the profile data (income level, occupation, education, caste category, disability status).
* **Saved Schemes (Wishlist)**: A visual panel to save schemes for quick access later, showing badge updates if a scheme's deadline is within 7 days.

### 2. Streamlined Multi-Lingual Interface (Hindi, Marathi, English)
* **Single-Tap Translation**: Smooth toggle across English, Hindi, and Marathi without page unmounts or state resets.
* **Semantic Parity**: Marathi and Hindi translations must maintain perfect semantic layout parity. Text bounds must not cause flex wrappers or navigation buttons to overflow or drop to separate rows.
* **Bilingual FAQ & Document Checklists**: Required documents and frequently asked questions must be fully localized to local languages.

### 3. Application Progress Tracker & Notifications
* **Visual Status Stepper**: Interactive step-by-step progress tracking for citizen applications:
  ```
  [Submitted] ───> [Under Verification] ───> [Decision: Approved / Rejected]
  ```
* **Offline Access (PWA)**: Support offline bookmark reading and cached scheme explorer search using Service Workers, enabling access in deep rural areas with spotty connectivity.
* **SMS & Push Alerts**: Automated SMS notifications to citizen phones when an application's status moves from `pending` to `approved` or `rejected`.

---

## 👑 Persona 2: The Administrator (Operations Requirements)
Government administrators require reliable tools to govern hundreds of schemes, categories, and review applications efficiently.

### 1. Unified Applications Management (Review Console)
* **Tactile Cohort Filters**: Ability to filter incoming applications by Status (All / Pending / Approved / Rejected) and Government Levels (Central / State).
* **Tactical Modal Inspections**: A detailed pop-up overlay modal showing the applicant's complete profile parameters (annual income, occupation, etc.) side-by-side with required documents.
* **Decision Remarks Prompt**: Force administrators to provide clear notes/remarks when rejecting an application, which is automatically broadcasted to the citizen's notification center.

### 2. Auto-Translation & Flexible Override
* **Auto-Translation Engine**: When adding or editing schemes, the system must translate English titles, benefits, documents, and FAQs into Hindi and Marathi via the Google Translate API.
* **Granular Manual Edits**: Admins must be able to overwrite translated Hindi/Marathi values in-place inside the modal inputs in case of grammar, dialect, or local idiom nuances.
* **Spinner Loading Shield**: Display a visually rich loading banner while translation fetches are active, shielding inputs from premature edits.

### 3. Cohort Notification Broadcast Engine
* **Segmented Broadcasts**: Send tailored system notifications to specific cohorts:
  - All Farmers (e.g., about PM Kisan upgrades)
  - All Students (e.g., about NSP scholarship deadlines)
  - Specific State Citizens (e.g., Maharashtra state schemes)
* **Audit Logs**: Maintain a list of all historical notification broadcasts for operational accountability.

---

## 💻 Persona 3: The CTO (Technical & Scalability Requirements)
For 100,000+ users, the backend must be built to support high-concurrency connections, secure sensitive citizen data, and remain highly resilient.

### 1. Database Optimizations & Query Safety
* **Indexes**: Build PostgreSQL indexes on frequently scanned columns:
  - `user_applied_schemes(user_id, status)`
  - `user_profiles(user_id)`
  - `schemes(category_id, is_active)`
* **Coalesce Scanning Protection**: Ensure Go SQL scanners never crash due to database NULL fields (e.g., using `COALESCE(a.notes, '')` for nullable text fields).
* **Connection Pooling**: Configure database connection pool parameters in `InitDB()` to scale up to 100 concurrent connections with timeouts.

### 2. Centralized API Architecture (Frontend Refactoring)
* **Modular Separation**: Decouple component/store layers from direct HTTP `fetch` requests. All API endpoints must be declared inside a specialized `frontend/src/api` module.
* **Universal HTTP Client**: Create a unified `client.js` HTTP wrapper that handles automatic JWT authorization injections, structured errors, and JSON marshaling.
* **Robust Error Handling**: Handle common network fallbacks, cross-origin resource sharing (CORS) rules, and API connection timeouts gracefully with UI toast alerts.

### 3. Production Security & Data Compliance
* **PII Encryption**: Encrypt sensitive citizen details (Aadhaar, income parameters) at-rest in the PostgreSQL database.
* **Token Rotation**: Implement HTTP-only cookie-based session management or JWT refresh tokens to avoid access token theft.
* **Rate Limiting**: Protect endpoints against brute force and DDoS attacks:
  - `/api/auth/login` (max 5 requests per minute per IP)
  - `/api/translate` (max 30 requests per minute per IP)
  - `/api/eligibility-check` (max 50 requests per minute per IP)

---

## 📊 Summary Gap Analysis: Current vs. Target

| Feature | Current State | Production Target Requirement | Action Status |
| :--- | :--- | :--- | :--- |
| **Admin Applications** | No console; Database query NULL bugs caused empty lists | Dedicated Applications Review console with status filters, remarks, and Null-safe queries | **Implemented (Fixed null bugs + Seeded mock data)** |
| **Card Navigation** | Static Overview numbers | Clickable stat cards routing to filtered datasets | **Implemented** |
| **API Decoupling** | API urls and `fetch` calls scattered inside Pinia store files | Centralized `src/api` module with structured modules and auto-headers | **Planned for execution** |
| **Translation Resiliency** | Relative path fetch URL bug returned 404; no load feedback | Correct absolute target, spinner indicator, and manual overwrite fields | **Implemented (Bug resolved + loading banner added)** |
| **Scale Enhancements** | Local SQLite/Postgres setups | Production indexing, rate-limiting, and PII encryption | **Added to Backlog** |
