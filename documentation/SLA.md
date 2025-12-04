## SLA – Whoknows

**Date:** 4 December 2025  
**Project:** Whoknows

---

### 1. What this service is

- **What we provide**: Access to the Whoknows application for our users.
- **Main features**: Login, registration, search, weather forecast, and the main page of the app.
- **Maintenance**: If we plan work that might break the app, we’ll try to let people know about **48 hours** before via email.

---

### 2. How well it should run

- **Uptime**: We aim for about **95.5%** uptime.
- **Speed**: Under normal internet, our site should load fast within **1 second**.
- **Backups**: We utilize Azure Backup Services for our database.

---

### 3. When something breaks

- **Big problems**:
  - We try to react within **1 hour** and aim to fix or at least reduce the impact within **4 hours**.
- **Medium problems** :
  - We try to react within **4 business hours** and aim to fix it within **1 business day**.
- **Small problems**:
  - We try to react within **1–3 business days** and fix them as part of our normal updates.

---

### 4. How we handle security

- We use logins and secure cookie sessions so only the right people can access their accounts.
- We use **HTTPS** so data between the user and the app is encrypted.
- We use Dependabot and Code Scanning in our repository to ensure dependencies are updated regularly.
- We use fail2ban and SSH to securely login to the servers.
- We tell users affected by data leak to change their passwords. Their profiles will not be blocked until the password is changed.

---

### 5. Rules and privacy

- We try to follow relevant privacy and data protection rules.
- When the project ends, we aim to delete or anonymize user data after 10th of January 2026.
