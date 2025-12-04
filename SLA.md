## DVK Application – Simple SLA

**Effective Date:** 4 december 2025
**Service:** DVK Web Application

---

### 1. Service Scope

- **Service**: Access to the DVK web app for registered users.
- **Features**: Login, registration, search, main pages, and related data storage.
- **Maintenance**: Planned maintenance with **48 hours** notice when possible.

---

### 2. Performance Metrics

- **Availability**: Target **99.5% uptime per month**, excluding planned maintenance and events outside our control.
- **Response Time**: Under normal load, most core requests aim to respond in **≤ 1 second** (within our infrastructure).
- **Backups**: Data backed up at least **daily** and retained for at least **7 days**.

---

### 3. Incident Response & Resolution

- **Critical** (full outage or core features unusable):
  - Acknowledgement within **1 hour**, target fix/mitigation within **4 hours**.
- **High** (major feature issues for many users):
  - Acknowledgement within **4 business hours**, target fix within **1 business day**.
- **Medium/Low** (minor or cosmetic issues):
  - Acknowledgement within **1–3 business days**, resolved in normal release cycles where feasible.

---

### 4. Security Measures

- Secure authentication and session management.
- HTTPS (TLS) for data in transit.
- Industry-standard protections for data at rest.
- Access to customer data limited to authorized personnel.

---

### 5. Compliance Standards

- Operated in line with applicable data protection laws (e.g., **GDPR**, **CCPA/CPRA**, where relevant).
- Security practices aligned with common industry standards (e.g., **ISO 27001**, **SOC 2**-style controls), where feasible.
