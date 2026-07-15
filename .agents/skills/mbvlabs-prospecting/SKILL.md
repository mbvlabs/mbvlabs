---
name: mbvlabs-prospecting
description: Finds and qualifies prospective MBV Labs consulting clients in the United States, Denmark, and England, identifies the best public business contact when possible, and saves an evidence-backed Markdown report. Use when asked to find MBV Labs prospects, prospective clients, target companies, qualified leads, founders or technical leaders to contact, or companies that may need fractional technical leadership or senior software engineering.
---

# MBV Labs Prospecting

Find a small number of companies that plausibly need MBV Labs now. A relevant, source-backed company without an email is more useful than a weak company with one.

This skill researches and writes a private Markdown report. It never sends email, submits forms, adds contacts to a CRM, or claims that a published address grants permission to contact someone.

## Load Context

Before searching, read completely:

1. `.agents/product-marketing.md`
2. `views/service_offerings_resource.templ`
3. `views/pages_resource.templ`

Use the current positioning as the source of truth. MBV Labs is a fractional tech lead who stays hands-on, with current service emphasis on backend systems, Go, APIs, integrations, internal platforms, AI workflows, technical leadership, and codebase modernization.

## Defaults

Unless the user provides different instructions:

- Find 15 qualified prospects.
- Split them evenly: 5 in the United States, 5 in Denmark, and 5 in England.
- Treat 5 to 60 employees as the ideal company size.
- Accept smaller companies only with clear budget, a real software need, and an accessible buyer.
- Accept companies with 61 to 150 employees only when offer fit and a current buying signal are strong.
- Usually exclude companies above 150 employees.
- Prefer quality over count. Return fewer prospects rather than padding the report.

Ask only for a missing choice that would materially change the run, such as a required sector, service focus, count, or market weighting. Otherwise use these defaults and continue.

## API Access

Use the MBV Labs API, not Serper or Firecrawl directly. All routes use HTTP Basic authentication.

Use `https://mbvlabs.com` unless `MBVLABS_API_BASE_URL` is set. Credentials come from `API_BASIC_AUTH_USERNAME` and `API_BASIC_AUTH_PASSWORD`. Never print, persist, or return credentials.

| Method | Route | Purpose |
|---|---|---|
| `POST` | `/api/search` | Discover companies, buying signals, people, and supporting sources |
| `POST` | `/api/scrape` | Read individual public company, team, careers, contact, job, and news pages |

Typical search requests:

```json
{"query":"startup hiring backend engineer product launch","location":"United States","countryCode":"us","languageCode":"en","num":10,"timeRange":"qdr:m"}
```

```json
{"query":"software startup søger udvikler AI produkt","location":"Denmark","countryCode":"dk","languageCode":"da","num":10,"timeRange":"qdr:y"}
```

```json
{"query":"founder led software company hiring platform engineer","location":"England","countryCode":"gb","languageCode":"en","num":10,"timeRange":"qdr:m"}
```

Typical scrape request:

```json
{"url":"https://example.com/about","formats":["markdown"],"onlyMainContent":true}
```

If the API or credentials are unavailable, report the exact missing prerequisite and stop. Do not invent prospects or silently call provider APIs directly.

## Ideal Client Profile

A company qualifies through the combination of team shape, relevant work, timing, and buyer access.

### Strong fit

- Founder-led software company, lean product team, or agency needing senior implementation capacity.
- A focused need involving backend systems, Go, APIs, integrations, data flows, internal tools, automation, AI workflows, product delivery, or modernization.
- Senior judgment and hands-on delivery matter more than adding a large delivery team.
- A visible current pressure such as a launch, funding event, new product, technical hiring, expansion, migration, operational bottleneck, or constrained engineering capacity.
- A founder, CEO, CTO, engineering lead, product lead, operations lead, or agency owner can plausibly buy the work.

### Disqualifiers

- More than 150 employees without an explicit user override.
- No evidence of a relevant software need or current business pressure.
- Pre-budget hobby project, equity-only request, cheapest-labour buyer, or no clear decision-maker.
- Large outsourced development company competing primarily on implementation labour.
- Inactive, closed, duplicated, outside the requested geography, or supported only by stale or unverifiable claims.

## Research Workflow

### 1. Build the candidate pool

Discover roughly twice the requested final count. Search several angles in each market rather than repeating one query:

- recent funding, launches, expansion, and product changes
- engineering, product, platform, AI, data, or operations hiring
- founder-led products and lean technical teams
- agencies needing backend, Go, AI, or senior delivery capacity
- migrations, integrations, automation, internal tooling, and modernization
- public descriptions of delivery pressure or manual operational work

Use both Danish and English queries for Denmark. For England, verify the company is in England and exclude Scotland, Wales, and Northern Ireland. Vary US searches enough to avoid returning only one city or startup hub.

Search results are discovery evidence, not final proof. Never qualify a company from a snippet alone.

### 2. Dedupe and triage

Dedupe by company domain before scraping. Remove obvious geography, size, service, and anti-persona misses first to save API credits.

For each promising candidate, identify the minimum pages needed to verify:

- official home or product page for what the company does
- official about or team page for ownership and approximate size
- official careers, changelog, blog, or job page for current pressure
- official contact page for a public business route
- one credible independent source when available for funding, launch, size, or timing

### 3. Verify with scraping

Scrape individual public pages only. Record what is observed, what is inferred, the source URL, the visible source date when available, and the date checked.

Estimate size as a band. Prefer official team and careers pages, credible company databases, and reputable reporting. If evidence is incomplete, use `Unknown` or a broad range with Low confidence. Never invent an exact headcount.

### 4. Map the best buyer

Choose one primary person, with an optional second person only when both have distinct buying roles:

- very small company: founder, CEO, or technical founder
- product company: CTO, head of engineering, engineering lead, or head of product
- workflow or internal-tool need: head of product, operations lead, or technical operator
- agency: owner, managing director, CTO, or technical director
- larger eligible company: the leader of the specific lean product or engineering team

Prefer names and roles published on the company's own site or in a public business announcement. Search public professional pages when necessary, but do not scrape restricted platforms.

Use only a business email published in a public professional context. Never guess an email pattern. Set contact status to one of:

- `Published business email`
- `Company contact route`
- `Not found`

Do not claim that an address is verified or deliverable. This project has no email verification API.

## Scoring

Score each candidate out of 100:

| Dimension | Points | Rules |
|---|---:|---|
| Offer fit | 30 | 25 to 30 for a direct core-service need; 15 to 24 for plausible adjacent work; below 15 is normally excluded |
| Size fit | 20 | 20 for 5 to 60; 14 for 61 to 100; 8 for fewer than 5 or 101 to 150; 5 when unknown; 0 and usually exclude above 150 |
| Why now | 25 | 20 to 25 for a specific recent signal; 10 to 19 for indirect or older pressure; below 10 lacks useful timing |
| Buyer reachability | 15 | 15 for a named buyer with published business email; 10 for a named buyer with a public route; 5 for a relevant role or generic company route; 0 for none |
| Evidence quality | 10 | 10 for official plus independent evidence; 7 for official evidence only; 4 for a credible third party only; 0 for unverifiable claims |

The primary shortlist requires 70 or more, no disqualifier, and at least Medium confidence. Missing email does not disqualify a strong company.

Confidence:

- `High`: official evidence plus credible corroboration
- `Medium`: one strong official or credible source with consistent supporting evidence
- `Low`: important facts remain ambiguous; keep outside the primary shortlist

## Compliance and Source Boundaries

Use public business data only. This is operational guidance, not legal advice.

- Never bulk scrape LinkedIn, Google Maps, Yelp, paywalled sources, login walls, or rate-limited platforms.
- Never bypass CAPTCHAs, bot protection, or access controls.
- Scraping an individual company's public website or a public news page is allowed.
- Never use personal email addresses, guessed addresses, breached data, unprovenanced lists, or sensitive personal traits.
- Capture the source URL and collection date for every named contact and email.
- For United States contacts, note that downstream sending must satisfy CAN-SPAM, sender identification, and opt-out requirements.
- For Denmark, never mark a contact cold-email ready merely because an address is public. Add `Consent or relationship review required`.
- For England, record whether the target appears to be a corporate subscriber or an individual subscriber such as a sole trader. Add `UK GDPR and PECR review required` when status or lawful basis is uncertain.
- Never describe a prospect as interested, ready to buy, or safe to email. Describe it as a potential client based on public evidence.

## Markdown Output

Write the final report with the `write` tool to:

```text
prospects/YYYY-MM-DD-mbvlabs-prospects.md
```

Add a short focus slug before `.md` when the run targets a specific service or sector. If the file already exists, add `-2`, `-3`, and so on rather than overwriting it. Keep reports under `prospects/` and treat public contact details as private repository data.

Use this structure:

```markdown
# MBV Labs prospect report

- Generated: YYYY-MM-DD
- Markets: United States, Denmark, England
- Requested: 15
- Qualified: 12
- Ideal size: 5 to 60 employees
- Focus: Default MBV Labs positioning

## ICP used

One concise paragraph describing the run-specific fit, buying signals, buyers, and disqualifiers.

## Top targets

| Score | Company | Market | Size | Why now | Best contact | Contact status | Confidence |
|---:|---|---|---|---|---|---|---|

## Prospect details

### 1. Company name, score/100

- Website:
- Location:
- Estimated size and confidence:
- Why MBV Labs fits:
- Why now:
- Best person and role:
- Contact: published email, company route, or Not found
- Contact source and date checked:
- Outreach review:
- Evidence: observed facts, followed by clearly labelled inference
- Sources: linked list with visible publication dates when available

## Relevant companies without a direct email

Include qualifying companies here when a useful person or company contact route exists.

## Search coverage and limitations

State candidate count, markets and query angles covered, major rejection reasons, API or source gaps, and facts that need manual confirmation.
```

Rank the full qualified list by score. In `Top targets`, show the best three to five with one specific reason each. A company may appear in both the compact top table and its detailed entry.

## Final Quality Check

Before writing the report, verify:

- Every shortlisted company has a source-backed MBV Labs fit and why-now signal.
- The ideal 5 to 60 band was prioritised without treating it as a hard boundary.
- No company above 150 was included without an explicit user override.
- England results are actually in England.
- Domains are deduplicated.
- Scores add up and match the evidence.
- Size is a sourced band with honest confidence.
- Every person and email has source lineage and a checked date.
- Missing emails remain `Not found` rather than guessed.
- Denmark and England include the correct outreach review note.
- Restricted platforms were not scraped.
- The report contains no credentials and no claim of consent, deliverability, or buying intent.
