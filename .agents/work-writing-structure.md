# Work Writing Structure

Use this guide when creating or rewriting entries for the `works` resource. Read `.agents/product-marketing.md` first so each Work supports the MBV Labs positioning: a fractional tech lead who combines senior engineering judgment with hands-on implementation.

## Goal

Each Work should show:

- The business or product context.
- The technical problem or ambiguity.
- What MBV Labs personally did.
- The engineering judgment behind the work.
- The practical outcome.

Write in first person by default: "I built", "I worked", "I shaped", "I helped ship". Keep the tone direct, calm, practical, and technically credible.

## Field Roles

### Title

Use a clear role-and-outcome title, not a vague project name.

Good patterns:

- `Principal engineering for [Client/Product]`
- `Backend systems for [Client/Product]`
- `[Capability] for [Client/Product]`
- `[Role] across [Product/Platform]`

Avoid:

- Cute campaign-style headlines.
- Generic labels like `AI Platform` or `Client Work`.
- Overclaiming beyond what was actually done.

### Summary

One sentence. State the role, client/product, main work, and strongest concrete outcome.

Formula:

`Worked as [role] on [client/product], building [system/work] that [specific outcome].`

Example:

`Worked as Principal Engineer on ChatSheet AI, building backend systems for multiple MVPs and enabling vector-based context search across several CMS platforms.`

### Challenge

This is the public `Overview` section. Keep it short and scannable.

Include:

- What the client/product was trying to do.
- Why the work was technically hard or ambiguous.
- What had to be true for the product to move forward.

Do not list everything delivered here. This section sets up the problem.

### Approach

Explain the engineering posture and decision-making.

Include:

- How I approached the problem.
- The tradeoffs that mattered.
- Any architecture, integration, data, AI, workflow, or delivery constraints.

Avoid repeating the deliverables list. This section should show judgment.

### Deliverables

Use a concise list-like sentence. Name the concrete work shipped or provided.

Include:

- Systems or features.
- Platforms or integrations.
- Technical leadership, architecture, implementation, audits, migrations, or tooling.

Keep this practical and specific.

### Outcome

Name the result without fabricating metrics.

Use:

- Shipped MVPs.
- Reduced manual work.
- Enabled a workflow.
- Improved maintainability.
- Supported launch, demos, sales, operations, or product iteration.
- Created reusable infrastructure or patterns.

If no metric is provided, describe the operational or product outcome plainly.

### Content

The long `content` field must not repeat `challenge`, `approach`, `deliverables`, and `outcome` in paragraph form. Use it to add depth.

Recommended structure:

1. `## [Role or central idea]`
   Describe the engagement context and why the work mattered.

2. `## [Specific technical theme]`
   Go deeper on one important system, workflow, or decision.

3. `## [Shipping or collaboration theme]`
   Explain how the work moved from ambiguity to shipped software.

4. `## [Tradeoff or lesson]`
   Show the engineering judgment behind the work.

Good content sections should:

- Tell a narrative, not restate metadata.
- Use concrete technical nouns.
- Explain tradeoffs in business terms.
- Make MBV Labs feel like a technical partner, not a ticket taker.
- Stay honest when exact metrics are unknown.

Avoid:

- Repeating the same phrases from the summary fields.
- Fabricated numbers, testimonials, or launch claims.
- Generic marketing language like "optimized", "streamlined", or "revolutionary".
- Turning the case study into implementation documentation.

## Required Input

Ask for missing details only when they materially change the quality or truth of the Work.

Useful questions:

- What should be the primary product/client name?
- What was the role or relationship?
- What dates should be shown?
- What were the concrete outcomes?
- What technologies or platforms mattered?
- Should any names, metrics, clients, or systems stay private?
- Are there cover or logo URLs?

Do not block on:

- Exact dates if month/year is enough.
- Metrics when a qualitative outcome is available.
- Cover image if the Work is staying draft.

## Draft Publishing Defaults

Unless explicitly told otherwise:

- `status`: `draft`
- `publishedAt`: empty
- `isFeatured`: `false`
- Leave unknown logo, cover, and date fields empty rather than inventing details.

## Quality Check

Before saving a Work, check:

- Does the summary explain the engagement in one sentence?
- Are the short fields scannable and non-repetitive?
- Does the long content add narrative depth?
- Is every outcome supported by provided facts?
- Does the copy show both technical leadership and hands-on execution?
- Is it written in the site's first-person voice?
