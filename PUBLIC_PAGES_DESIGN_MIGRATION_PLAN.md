# Public Pages Cobalt Design Migration Plan

Reference commit: `2eebf702ccd4`

## Goal

Move every public server-rendered HTML page into the cobalt field-report design established by `views/landing.templ`.

Implement and review one view at a time. Preserve all page behavior, metadata, dynamic content, and form interactions. Do not change the Inertia admin interface.

## Current State

The landing redesign already created the migration seam:

- `views/landing.templ` uses `themedBase("cobalt-grid", ...)`.
- `.cobalt-grid` scopes the public palette and typography.
- The shared public header and footer use semantic theme tokens.
- Every other public Templ view still uses `base(...)` and hard-coded white, zinc, or gray styling.
- Admin is a separate Inertia/Vue surface under `resources/js` and uses the root light/dark tokens.

```text
PUBLIC HTML                                     ADMIN
controllers/*.go                                 controllers/admin/*.go
       |                                                   |
       v                                                   v
views/*.templ  -> themedBase("cobalt-grid")       resources/js/Pages/Admin/*.vue
       |                                                   |
       +---- scoped .cobalt-grid tokens                    +---- root/.dark tokens
```

## Scope

“All public pages” means every server-rendered Templ HTML state reachable outside `/admin`.

### Included

- About
- Services
- Contact
- Work index and show
- Project index and show
- Writing index and show
- Login
- Password request and reset
- Email confirmation
- Not found
- Bad request
- Internal error

The landing page is already complete and remains the reference implementation.

Account access screens are included because they are public Templ routes. Only their presentation changes. Session, password reset, confirmation, redirect, and admin behavior remain untouched.

### Excluded

- All `/admin` Inertia pages and components
- Admin controllers and routes
- APIs
- Sitemap and robots responses
- Static assets as separate page designs
- Email templates
- Legacy redirects as separate designs
- Models, migrations, and database behavior

## Admin Isolation Rules

- Do not change files under `resources/js` for this migration.
- Do not change `controllers/admin` or `router/routes/admin*.go`.
- Do not make cobalt the root theme.
- Do not alter generic `:root`, `.light`, or `.dark` behavior.
- Keep new custom CSS under `.cobalt-grid` or use uniquely named `cobalt-*` classes.
- Add no global element overrides that could alter admin rendering.
- Smoke-check `/admin` and one admin resource after each public view.

## Design Contract

Apply the landing page's visual grammar without forcing every page into the exact landing layout.

### Visual Language

- Cream background and cobalt foreground from `.cobalt-grid`
- Newsreader for display headings
- Hanken Grotesk for body copy
- DM Mono for metadata
- Large editorial typography
- Uppercase field labels
- Numbered rows and ledger metadata
- One-pixel rules
- Square corners
- No card shadows
- Restrained hover fills
- Wide `container` layout with `px-6 md:px-8`

### Implementation Rules

- Switch only the target view from `base(...)` to `themedBase("cobalt-grid", ...)`.
- Replace hard-coded zinc, gray, and white classes with semantic theme utilities.
- Reuse `cobalt-display`, `cobalt-section-title`, and `cobalt-row-title`.
- Add a scoped `cobalt-*` CSS helper only after a real repeated need appears.
- Do not create a generic component library before two migrated views prove the repetition.
- Preserve existing page fragment IDs and container IDs.
- Preserve all head options and schema builders.

Example:

```templ
templ (wi WorkIndex) Page() {
    @themedBase(
        "cobalt-grid",
        // Existing title, description, slug, and schema options stay unchanged.
    ) {
        @templ.Fragment(wi.PageFragment()) {
            <main class="overflow-hidden bg-background text-foreground">
                <p class="font-mono text-xs uppercase tracking-[0.18em]">...</p>
                <h1 class="cobalt-display">...</h1>
            </main>
        }
    }
}
```

## Migration Queue

Each numbered item is a separate implementation and review gate.

### 01. Work Index ✅

- Route: `GET /work`
- Source: `views/works_resource.templ`
- Use the landing page's selected-work ledger as the starting pattern.
- Convert the hero, client rows, people list, empty state, and experience history into field-report sections.
- Preserve item loops, dates, logos, outcomes, links, and long experience content.

### 02. Work Show ✅

- Route: `GET /work/:slug`
- Source: `views/works_resource.templ`
- Establish the detail-page grammar: editorial hero, three-column metadata ledger, media, side rail, and narrative sections.
- Preserve optional client data, dates, tags, challenge, approach, deliverables, outcome, and Markdown content.

### 03. Project Index ✅

- Route: `GET /projects`
- Source: `views/projects_resource.templ`
- Adapt the work ledger to project type, source, year, technologies, logo or initial, and live-site actions.
- Preserve the empty collection state and conditional external links.

### 04. Project Show ✅

- Route: `GET /projects/:slug`
- Source: `views/projects_resource.templ`
- Reuse the detail grammar established by Work Show without extracting speculative components.
- Preserve Markdown, technologies, source and live availability, image fallback, dates, and external links.

### 05. Writing Index

- Route: `GET /blog`
- Source: `views/blog_posts_resource.templ`
- Present posts as dated dispatch entries with mono topic metadata and editorial titles.
- Preserve tag parsing, excerpts, publication dates, article links, and the empty state.

### 06. Writing Show

- Route: `GET /blog/:slug`
- Source: `views/blog_posts_resource.templ`
- Build a readable cobalt article layout with restrained measure and scoped prose styling.
- Preserve cover image, published and updated rules, topics, Markdown output, schema, and back navigation.
- Validate headings, links, code, preformatted blocks, lists, blockquotes, and images in rendered Markdown.

### 07. Services

- Route: `GET /services`
- Source: `views/service_offerings_resource.templ`
- Turn offerings into a numbered service ledger with fit, process, background, and CTA field sections.
- Reuse existing content and route actions.
- Do not invent new offer data or controller structure.

### 08. About

- Route: `GET /about`
- Source: `views/pages_resource.templ`
- Compose a personal field report around what I do, audience, working style, and founder quote.
- Keep title, description, breadcrumbs, and contact route unchanged.

### 09. Contact

- Routes: `GET /contact` and `POST /contact`
- Source: `views/project_inquiries_resource.templ`
- Style the inquiry as a structured intake sheet with square controls, visible labels, strong focus states, and a clear submit action.
- Preserve field names, action, method, required/min/max constraints, autocomplete, options, and message behavior exactly.
- Check keyboard order, invalid fields, submitting state, success flash, and narrow layouts.

### 10. Login

- Route: `GET /users/sign_in`
- Source: `views/login.templ`
- Replace the gray rounded card with a compact cobalt access form.
- Preserve the Datastar submit action, bindings, disabled state, password link, noindex metadata, and post-login redirect behavior.

### 11. Password Request

- Route: `GET /users/password/new`
- Source: `views/reset_password.templ`
- Apply the access-form grammar to the email request state.
- Preserve submit action, binding, disabled state, login link, and noindex metadata.

### 12. Password Reset

- Route: `GET /users/password/:token/edit`
- Source: `views/reset_password.templ`
- Migrate the new-password state independently from the request view, even though both share one source file.
- Preserve the hidden token binding, both password fields, PUT action, disabled state, and noindex metadata.

### 13. Email Confirmation

- Route: `GET /users/confirmation/new`
- Source: `views/confirm_email.templ`
- Use the access-form grammar for the six-digit code.
- Preserve maxlength, binding, POST action, loading state, and noindex metadata.

### 14. Not Found

- Route: unmatched public route
- Source: `views/not_found.templ`
- Create a concise branded field-note error with a route back home and useful navigation.
- Ensure it works for any unmatched public path.

### 15. Bad Request

- State: 400 render path
- Source: `views/bad_request.templ`
- Use the system-page grammar with clear recovery copy.
- Do not expose submitted data, parser details, or internal errors.

### 16. Internal Error

- State: 500 render path
- Source: `views/internal_error.templ`
- Create a calm branded failure state with a safe recovery action.
- Do not leak stack traces, storage errors, or request internals.

## Per-View Workflow

```text
READ TARGET + ALL DATA BRANCHES
          |
          v
COMPOSE WITH EXISTING COBALT TOKENS
          |
          v
GENERATE TEMPL + CSS
          |
          v
CHECK DATA, RESPONSIVE, KEYBOARD, ADMIN BOUNDARY
          |
          v
REVIEW GATE -------------- approved? --------------> NEXT VIEW
          |
          +---------------- no: revise this view only
```

Stop after each view for review before starting the next one.

## Per-View Acceptance Criteria

- The target view uses `themedBase("cobalt-grid", ...)`.
- No old zinc, gray, hard-coded white, rounded-card, or shadow styling remains inside the migrated view.
- Existing page fragment ID and container ID remain stable.
- Existing title, description, slug, image, metadata, and schema options remain intact.
- Dynamic branches render correctly.
- Links, forms, and Datastar attributes remain functional.
- The page has one H1 and a logical heading order.
- Controls have labels and visible keyboard focus.
- The layout works at 375px, 768px, and 1440px.
- Header, mobile navigation, footer, and flash messages remain usable.
- Admin has no visual or behavioral change.

## Content States to Review

- Empty and populated indexes
- Detail records with and without images
- Missing and present dates
- Missing and present tags
- Missing and present optional narrative sections
- Missing and present external links
- Long titles and metadata
- Long body content
- Rendered Markdown
- Form idle, keyboard focus, invalid, submitting, and server flash states

## Verification After Every View

```bash
./bin/templ generate
./bin/tailwindcli -i css/base.css -o assets/css/style.css
go vet ./...
andurel doctor --json
```

Run `gofmt` only on changed non-generated Go files, if any.

### Browser Smoke Check

- Open the target route at 375px, 768px, and 1440px.
- Exercise target links, optional content, and form interactions.
- Confirm the shared header, mobile navigation, footer, and flash display.
- Open `/admin` and one admin resource to confirm no change.
- Inspect the diff for the target source, generated Templ Go, and compiled CSS only.

## Expected File Boundary

### Primary Sources

- `views/works_resource.templ`
- `views/projects_resource.templ`
- `views/blog_posts_resource.templ`
- `views/service_offerings_resource.templ`
- `views/pages_resource.templ`
- `views/project_inquiries_resource.templ`
- `views/login.templ`
- `views/reset_password.templ`
- `views/confirm_email.templ`
- `views/not_found.templ`
- `views/bad_request.templ`
- `views/internal_error.templ`

### Conditional Shared Sources

Change these only when the target view proves the need:

- `css/custom.css` for scoped cobalt helpers
- `views/layout.templ` if a shared public-shell defect is discovered
- `css/themes.css` if the cobalt theme lacks a required semantic token

Keep `views/landing.templ` unchanged as the reference implementation.

### Generated Outputs

- Matching `views/*_templ.go` files
- `assets/css/style.css`

## Starting Point

Begin with **Stage 01: Work Index**.

It reuses the strongest existing landing-page pattern and proves the migration workflow with the least new design code.
