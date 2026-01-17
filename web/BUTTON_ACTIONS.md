# Button Actions Reference

## Hero Section (`/web/components/Hero.tsx`)

### 1. **Get Started Button**

- **Action**: Smooth scroll to the installation section
- **Implementation**: `onClick={() => document.getElementById('docs')?.scrollIntoView({ behavior: 'smooth' })}`
- **Target**: Scrolls to the `#docs` section in `GettingStarted.tsx`

### 2. **Audit on GitHub Button**

- **Action**: Opens GitHub repository in new tab
- **Link**: `https://github.com/PrakarshSingh5/FintechKit`
- **Implementation**: `<a>` tag with `target="_blank"` and `rel="noopener noreferrer"`

## Getting Started Section (`/web/components/GettingStarted.tsx`)

### 3. **Read the Full Documentation Button**

- **Action**: Opens GitHub README in new tab
- **Link**: `https://github.com/PrakarshSingh5/FintechKit#readme`
- **Implementation**: `<a>` tag with `target="_blank"` and `rel="noopener noreferrer"`

### 4. **Copy Install Command Button**

- **Action**: Copies `go get github.com/fintechkit/fintechkit` to clipboard
- **Implementation**: Uses `navigator.clipboard.writeText()` with visual feedback

---

## Alternative Button Actions You Could Use

### Option A: Link to Specific Documentation Pages

```tsx
<a href="https://github.com/PrakarshSingh5/FintechKit/blob/main/QUICK_REFERENCE.md">
  Quick Start Guide
</a>
```

### Option B: Link to Examples

```tsx
<a href="https://github.com/PrakarshSingh5/FintechKit/tree/main/examples">
  View Examples
</a>
```

### Option C: Link to Specific Integration Guide

```tsx
<a href="https://github.com/PrakarshSingh5/FintechKit/blob/main/RAZORPAY_INTEGRATION_GUIDE.md">
  Razorpay Integration
</a>
```

### Option D: Create a Dedicated Docs Page

If you want a full documentation site, you could:

1. Create a `/docs` route in your Next.js app
2. Link the button to `/docs`
3. Build out comprehensive documentation pages

Example:

```tsx
<Link href="/docs/getting-started">Get Started</Link>
```

---

## Current Flow

1. User lands on homepage
2. Clicks "Get Started" → Smoothly scrolls to installation section
3. Sees `go get` command → Can copy it
4. Clicks "Read the Full Documentation" → Opens GitHub README
5. Clicks "Audit on GitHub" → Opens repository to review code

This creates a smooth onboarding experience! 🚀
