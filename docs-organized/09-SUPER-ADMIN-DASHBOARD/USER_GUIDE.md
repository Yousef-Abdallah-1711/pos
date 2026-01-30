# Super Admin Dashboard - User Guide

This guide explains how to use the POS SaaS Super Admin Dashboard to manage your platform.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Navigation](#navigation)
3. [Dashboard Overview](#dashboard-overview)
4. [Managing Tenants](#managing-tenants)
5. [Managing Users](#managing-users)
6. [Impersonation Feature](#impersonation-feature)
7. [Audit Logs](#audit-logs)
8. [Settings](#settings)
9. [Security Tips](#security-tips)

---

## Getting Started

### Accessing the Dashboard

1. Open your browser and navigate to: `http://localhost:3001`
2. You will be redirected to the login page
3. Enter your super admin credentials:
   - **Email**: superadmin@pos-saas.com
   - **Password**: SuperAdmin@123

### First Login

After your first login:
1. You'll see the main dashboard with platform statistics
2. Configure your IP whitelist in Settings (recommended for production)
3. Change your password immediately
4. Review recent audit logs

---

## Navigation

### Sidebar Menu

The left sidebar contains the main navigation:

| Icon | Menu Item | Purpose |
|------|-----------|---------|
| 📊 | **Dashboard** | View platform statistics and recent activity |
| 🏢 | **Tenants** | Manage all customer accounts |
| 👥 | **Users** | View and manage all platform users |
| 📋 | **Audit Logs** | View complete activity history |
| ⚙️ | **Settings** | Configure account settings and IP whitelist |
| 🚪 | **Logout** | Sign out of the dashboard |

---

## Dashboard Overview

The main dashboard displays key metrics:

### Statistics Cards

**Total Tenants**
- Number of active customer accounts
- Click to view full tenant list

**Total Users**
- Count of all users across all tenants
- Includes users from all accounts

**Active Impersonations**
- Currently active impersonation sessions
- Shows real-time activity

**Total Audit Logs**
- Complete record of all super admin actions
- Used for compliance and troubleshooting

### Quick Actions

- **View All Tenants** - Access tenant management
- **View All Users** - Manage all platform users
- **View Audit Logs** - Review activity history
- **Manage Settings** - Configure account

### Recent Activity

Shows the 5 most recent actions performed:
- Login/logout events
- Impersonation sessions
- Tenant changes
- Settings modifications
- IP whitelist updates

---

## Managing Tenants

### Viewing Tenants

1. Click **Tenants** in the sidebar
2. You'll see a list of all customer accounts
3. Each tenant shows:
   - Business name
   - Unique slug identifier
   - Account creation date

### Searching Tenants

Use the search box to find tenants by:
- Name (business name)
- Slug (identifier)

Example: Type "acme" to find ACME Corporation

### Viewing Tenant Details

1. Click **View** on any tenant row
2. You'll see:
   - Full tenant information
   - Creation and update dates
   - List of all users in the tenant

### Tenant Users

On the tenant detail page:
- See all users in that specific tenant
- View their email, role, and creation date
- **Impersonate users** directly (see Impersonation section)
- Pagination for large teams

### Tenant Management

Currently available from the detail page:
- View all information
- See user metrics
- Impersonate team members

---

## Managing Users

### Viewing All Users

1. Click **Users** in the sidebar
2. See a complete list of all platform users
3. Users are shown with:
   - Full name (First + Last)
   - Email address
   - Role (Admin, User, etc.)
   - Tenant ID they belong to
   - Account creation date

### Searching Users

Use the search box to find users by:
- Name (first or last name)
- Email address

Example: Type "john@acme.com" or "John"

### User Details

Click **View** on any user to see:
- Full user information
- Email and role
- Tenant assignment
- Account creation date
- Impersonation option

---

## Impersonation Feature

The impersonation feature allows you to access any tenant's dashboard as if you were that user. This is helpful for:
- Technical support and troubleshooting
- Testing workflows
- Assisting customers with issues
- Administrative setup

### Starting an Impersonation Session

**Method 1: From Tenant Detail Page**
1. Navigate to **Tenants**
2. Click **View** on a tenant
3. Find the user you want to impersonate
4. Click **Impersonate** button

**Method 2: From Users List**
1. Navigate to **Users**
2. Find the user
3. Click **View**
4. Click **Impersonate** button

### Impersonation Modal

When you click Impersonate:
1. A confirmation window appears
2. Shows the target user's information
3. Displays a warning about full access
4. Click **Start Impersonation**

### Using the Impersonation Session

After starting impersonation:
1. A unique token is generated
2. Dashboard URL is displayed
3. Click **Open Dashboard in New Tab** to access the tenant dashboard

### Impersonation Banner

While impersonating:
- A yellow banner appears at the top
- Shows: "Impersonating User" with their name and email
- Displays session expiry time
- **End Session** button to terminate impersonation

### Session Duration

- Impersonation sessions last **1 hour**
- Session expires automatically after 1 hour
- You can manually end the session anytime
- Impersonation is fully logged in audit logs

### Important Notes

- ⚠️ You have **full access** while impersonating
- All actions are logged with your super admin ID
- Cannot impersonate other super admins
- Other users will not see you as the original user
- Session is specific to this super admin account

---

## Audit Logs

### What are Audit Logs?

Audit logs are a complete, immutable record of all super admin actions. They're essential for:
- Compliance and regulatory requirements
- Security audits
- Troubleshooting issues
- Detecting suspicious activity

### Accessing Audit Logs

1. Click **Audit Logs** in the sidebar
2. See all recorded actions
3. Most recent actions appear first

### Log Details

Each log entry shows:
- **Action**: Type of action (LOGIN_SUCCESS, IMPERSONATE_START, etc.)
- **Timestamp**: When it happened
- **IP Address**: Source IP of the request
- **Target User**: User affected (if applicable)
- **Target Tenant**: Tenant affected (if applicable)
- **Details**: Additional metadata

### Filtering Logs

Use the **Filter by Action** dropdown:
- **LOGIN_SUCCESS** - Successful logins
- **LOGIN_FAILED** - Failed login attempts
- **LOGOUT** - Logout events
- **IMPERSONATE_START** - Impersonation started
- **IMPERSONATE_END** - Impersonation ended
- **TENANT_DISABLED** - Tenant disabled
- **TENANT_ENABLED** - Tenant re-enabled
- **IP_WHITELIST_UPDATED** - IP settings changed

### Log Colors

Color-coded status badges:
- 🟢 **Green**: Successful actions (login, enable)
- 🔴 **Red**: Failed or negative actions (failed login, disable)
- 🟠 **Orange**: Impersonation start
- 🟣 **Purple**: Impersonation end
- 🟡 **Yellow**: Configuration changes

### Pagination

Navigate through pages of logs:
- Shows page number and total pages
- Each page has 50 entries
- Use Previous/Next buttons

---

## Settings

### Account Information

View your super admin account details:
- **Email**: Your login email
- **Account Status**: Active or Inactive

### IP Whitelist

Control which IP addresses can access your super admin account.

#### Adding IPs

1. Enter an IP address in the form:
   - **Single IP**: `192.168.1.1`
   - **CIDR Range**: `10.0.0.0/8`
2. Click **Add IP**
3. IP appears in the whitelist list
4. Click **Save Changes**

#### Removing IPs

1. Find the IP in the list
2. Click the **Delete** (trash) icon
3. IP is removed immediately
4. Click **Save Changes**

#### Empty Whitelist

- If the list is empty, **all IPs are allowed**
- This is useful during initial setup
- Add IPs for production environments

### Saving Changes

Click **Save Changes** button to apply all modifications:
- New IPs are added
- Removed IPs are deleted
- Changes are logged in audit logs

### Security Best Practices

The dashboard displays recommendations:
- ✓ Always enable IP whitelist for production
- ✓ Use static IP addresses when possible
- ✓ Monitor audit logs regularly
- ✓ Change password frequently
- ✓ Never share login credentials
- ✓ Log out when finished

---

## Security Tips

### Best Practices

**Account Security**
1. Change default password immediately
2. Use a strong, unique password
3. Never share your credentials
4. Enable IP whitelist in production

**Session Security**
1. Always log out when finished
2. Close the browser tab if leaving the dashboard
3. Don't leave your computer unattended
4. Use secure, private internet connections

**Impersonation Safety**
1. Only impersonate for legitimate business reasons
2. Document why you're impersonating
3. End sessions promptly after use
4. Review impersonation audit logs regularly

**Account Monitoring**
1. Review audit logs daily in production
2. Check for unusual login patterns
3. Monitor failed login attempts
4. Track IP whitelist changes

### Security Features

The platform includes built-in protection:
- 🔒 JWT tokens expire after 15 minutes
- 🔐 Failed logins lock account for 30 minutes
- 📊 All actions are logged immutably
- 🌐 IP whitelist validation
- ⏱️ Rate limiting (30 requests/minute)
- 🔑 Separate encryption for impersonation tokens

### Incident Response

If suspicious activity detected:

1. **Immediately log out** using the Logout button
2. **Review audit logs** for unauthorized actions
3. **Change your password** on the login page
4. **Contact your system administrator**
5. **Check IP whitelist** for unauthorized IPs
6. **Monitor impersonation sessions** for misuse

---

## Troubleshooting

### Login Issues

**"Invalid credentials" error**
- Verify email and password are correct
- Remember: case-sensitive password
- Reset password if needed

**"Account is locked" error**
- Account locked after 3 failed attempts
- Wait 30 minutes before retry
- Or contact administrator

**Cannot access dashboard**
- Check if your IP is whitelisted
- Verify you have internet connection
- Try a different browser or private window

### Impersonation Issues

**Impersonation token expired**
- Impersonation sessions expire after 1 hour
- Start a new impersonation session
- Save any work before session expires

**Cannot impersonate user**
- Verify user exists and is active
- Cannot impersonate other super admins
- Check audit logs for details

### Performance Issues

**Dashboard loads slowly**
- Check your internet connection
- Try refreshing the page
- Clear browser cache
- Try a different browser

**Search not working**
- Make sure search term is not empty
- Try more specific search terms
- Use exact email or partial name

---

## FAQ

### Q: How long do sessions last?
**A**: Authentication tokens last 15 minutes. Impersonation sessions last 1 hour.

### Q: What if I forget my password?
**A**: Contact your system administrator for password reset.

### Q: Can other users see I'm impersonating them?
**A**: No, but all impersonation is logged in audit logs. The target user will not receive notifications.

### Q: What information is logged in audit logs?
**A**: All super admin actions including logins, logouts, impersonations, tenant changes, and settings updates.

### Q: Can I impersonate another super admin?
**A**: No, impersonation is restricted to regular users only for security reasons.

### Q: What happens if my IP is not whitelisted?
**A**: Your login request will be rejected with an access denied error.

### Q: How do I export audit logs?
**A**: Currently, audit logs are viewed in the dashboard. Pagination supports viewing all logs. Contact your administrator for export options.

---

## Support

For technical support or questions:
1. Check the audit logs for error details
2. Review this user guide
3. Contact your system administrator
4. Check system status page for incidents

---

**Last Updated**: 2024-01-30
**Dashboard Version**: 1.0.0
**Status**: Production Ready
