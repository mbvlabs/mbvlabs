---
name: mbvlabs-article
description: Researches target keywords, outlines, drafts, and saves search-focused MBV Labs articles from supplied keywords, a general idea, or a rough outline. Use when creating blog posts intended to attract qualified consulting traffic, finding keywords from an initial prompt, conducting keyword or competitor research for MBV Labs, matching Morten's writing voice, or saving an approved article draft through the MBV Labs API.
---

# MBV Labs Article

Create useful technical articles that attract qualified founders, product leaders, and engineering leads. Traffic without relevance to MBV Labs services is not success.

This is a manual workflow with two mandatory approval gates:

1. Research and outline approval before drafting.
2. Draft approval before saving through the API.

Never skip or combine these gates, even if the initial request asks for an article. Never publish an article. The API creates drafts only.

Keyword targeting is mandatory. Every article must have one approved primary target keyword and a small supporting keyword set before outlining. Use keywords supplied by the user as the starting set. When none are supplied, derive candidates from the initial idea, question, argument, or outline and validate them with live search research.

## Load Context

Before research:

1. Read `.agents/skills/copywriting/SKILL.md` completely and apply it throughout this workflow.
2. Read `.agents/product-marketing.md` completely.
3. Read `references/voice.md` completely.
4. Read `references/structure.md` completely.
5. Inspect existing MBV Labs content relevant to the topic before choosing an angle.

Treat the marketing context as the source of truth for the offer and audience. `references/voice.md` is the complete, static prose voice reference. Diary entries provide current opinions and raw vocabulary, not finished prose. MBV Labs pages define positioning, not personal voice.

## API Access

Use the MBV Labs API, not Serper or Firecrawl directly. All API routes use HTTP Basic authentication.

Use `https://mbvlabs.com` by default. If `MBVLABS_API_BASE_URL` is set, use it instead. Credentials come from `API_BASIC_AUTH_USERNAME` and `API_BASIC_AUTH_PASSWORD`. Never print, persist, or include credentials in an answer.

Available routes:

| Method | Route | Purpose |
|---|---|---|
| `POST` | `/api/search` | Google results through Serper |
| `POST` | `/api/map` | Discover relevant pages on a domain through Firecrawl |
| `POST` | `/api/scrape` | Extract a page as Markdown through Firecrawl |
| `GET` | `/api/diary/thoughts/current-week` | Current first-person voice material |
| `POST` | `/api/articles` | Save an approved article as a draft |

Typical requests:

```json
{"query":"fractional tech lead for startups","languageCode":"en","num":10}
```

```json
{"url":"https://mbvlabs.com","search":"fractional tech lead","limit":50}
```

```json
{"url":"https://example.com/article","formats":["markdown"],"onlyMainContent":true}
```

If the API or credentials are unavailable, report the exact missing prerequisite. Do not replace live research with invented results or silently call provider APIs directly.

## Input

Accept a general idea, question, argument, or rough outline. Ask only for missing information that materially affects truth or direction, such as:

- a real experience the article is expected to describe
- the intended buyer when more than one is plausible
- a target geography when the query is location-dependent
- confidential details or claims that must not be used

Treat any supplied keywords as first-class requirements. Preserve their wording, research their intent, and account for each one in the research brief as the primary target, a supporting keyword, or a rejected keyword with a reason. Do not silently replace a supplied keyword because another phrase seems easier to write for.

When no keywords are supplied, extract seed terms from the initial prompt in the likely reader's language. Research those terms and their problem-aware, commercial, and long-tail variants. Do not ask the user for a keyword when live research can identify the best target.

Default to English, an international audience, and organic search traffic.

## Phase 1: Research

### 1. Check site fit

Map and inspect relevant pages on `mbvlabs.com`, including existing blog posts, services, work, and projects.

Determine:

- which MBV Labs service or proof point the topic supports
- the likely reader and their business pressure
- existing pages that could compete with or support the article
- natural internal links
- whether the idea is too broad, unrelated, or duplicative

If the idea has little connection to a plausible client problem, say so and propose a closer angle before spending more API credits.

### 2. Set the voice plan

Use only `references/voice.md` and material approved for the current task. Do not look up old articles.

Choose the closest voice shape and define:

- how the article will open
- which supplied experience or opinion gives it first-person authority
- where technical detail will connect to a practical consequence
- the strongest counterargument
- the boundary or uncertainty Morten should acknowledge
- the appropriate level of humor
- how the conclusion will state a current practical position

Do not reduce this step to adjectives such as `direct` or `practical`. If the planned article needs a personal fact, opinion, client example, or result that has not been supplied, ask Morten before drafting.

### 3. Build and validate the keyword set

Start from the keyword source:

- **Supplied keywords:** search each supplied keyword. Keep the original wording visible throughout the research brief.
- **No supplied keywords:** turn the initial prompt into several candidate phrases that a likely buyer would search. Include the topic term, a problem-aware phrase, a decision or commercial phrase, and a specific long-tail question.

Use organic results, People Also Ask questions, and related searches to identify the language searchers and competing pages use. Search only enough variants to compare intent and understand the result pages.

For every candidate primary keyword, assess:

- whether the dominant intent matches the proposed article
- whether the likely searcher resembles an MBV Labs buyer
- whether MBV Labs can credibly provide a better or more useful answer
- whether existing MBV Labs content already targets it
- which supporting keywords and questions belong to the same article

The search API does not provide keyword volume or keyword difficulty. Never invent either. Describe opportunity using observable signals such as intent match, result quality, repeated subtopics, weak coverage, and relevance to likely buyers.

### 4. Inspect competitors

Choose three to five genuine editorial competitors from the live results. Prefer pages that satisfy the same intent. Exclude irrelevant homepages, forums, social posts, and duplicate domains unless they dominate the result page for a useful reason.

Scrape each selected page. Record:

- title and URL
- apparent search intent
- H2 and H3 coverage
- useful examples, evidence, tables, checklists, or code
- where the page is vague, dated, overly developer-focused, or disconnected from business decisions
- claims that need independent verification

Competitor content is a map of reader expectations, not source material to rewrite. Do not copy phrasing, examples, or structure section by section.

### 5. Find a defensible angle

Select exactly one primary target keyword and a small set of supporting keywords and questions. The primary keyword determines the search intent, title direction, opening answer, and core outline. Prefer the overlap of:

- clear search intent
- a problem MBV Labs can credibly solve
- a useful angle supported by Morten's experience or judgment
- gaps in current results
- a natural path to a service, work example, project, or contact page

Qualified relevance beats estimated reach. A narrow query used by potential clients is often better than a broad query used mainly by students or developers looking for snippets.

If a supplied keyword does not fit the article, explain the intent mismatch or cannibalization risk and recommend where it belongs instead. Do not discard it silently. Do not proceed to the outline until the brief contains a defensible primary target keyword.

For factual claims, prefer primary sources, official documentation, or original research. Keep a source URL beside every planned statistic, dated fact, or attributed claim. If a claim cannot be verified, remove it or clearly frame it as opinion.

## Research Brief and Outline

Present the following before writing any prose draft:

### Research brief

- **Working idea:** one sentence
- **Target reader:** role, situation, and level of technical knowledge
- **Client relevance:** how this can lead to qualified MBV Labs work
- **Search intent:** what the reader wants to understand, decide, or do
- **Keyword source:** supplied by the user or discovered from the initial prompt, including the original supplied or seed terms
- **Primary target keyword:** one exact phrase
- **Supporting keywords and questions:** concise list drawn from the input and live results
- **Keyword decision:** why the primary target won, how each supplied keyword was assigned, and which candidates were rejected
- **SERP evidence:** key People Also Ask, related-search, and result patterns supporting the keyword decision
- **Competitor findings:** compact table with page, strengths, and gap
- **Existing site fit:** cannibalization risk, supporting pages, and internal-link opportunities
- **Distinct angle:** the useful claim or perspective competitors do not cover well
- **Voice plan:** opening, approved first-person source, technical consequence, counterargument, uncertainty boundary, humor level, and conclusion style
- **Structure choice:** selected pattern from `references/structure.md` and why it fits this intent
- **Evidence needed:** facts, examples, screenshots, code, or first-person details still required

State explicitly that traffic volume and keyword difficulty are unknown unless a real source supplied them.

### Proposed article package

- one recommended title and two alternatives
- proposed slug
- one-sentence promise to the reader
- article type, such as guide, decision framework, comparison, technical explainer, or opinion backed by experience
- primary conversion action and why it is natural

### Detailed outline

List every planned H2 and H3. Under each heading include:

- the question or job the section resolves
- the key point or argument
- concrete evidence, example, tradeoff, or source to use
- relevant supporting keyword or question when applicable
- natural internal link when applicable

Then stop and ask for outline approval or requested changes. Do not write the introduction, sample sections, or full article yet.

## Approval Gate 1

Proceed only after the user clearly approves the outline. If they request changes, revise the research brief or outline and ask again.

Outline approval permits drafting. It does not permit saving.

## Phase 2: Draft

### Structure

Follow the approved structure selected from `references/structure.md`. That reference captures what is useful in Outrank's examples, including early answers, headings based on reader jobs, complete decision coverage, information-dense tables or checklists, and selective FAQs. It also rejects Outrank's common padding, repetitive headings, generic analogies, unsupported statistics, and oversized promotional sections.

The article type determines the structure. Do not flatten an opinion into a generic how-to guide or inflate a focused answer into a 4,000-word article. Every H2 must perform a different reader job. Remove a section when deleting it would not make the reader's task harder.

### MBV Labs voice

Follow `references/voice.md` and the voice plan approved in Phase 1. The target is not a collection of brand adjectives. Reproduce the documented reasoning pattern:

1. Begin from a real experience, previous view, project, or decision when the evidence supports one.
2. Make the current tension or position clear early.
3. Use concrete technical detail to explain a consequence for understanding, ownership, delivery, reliability, cost, or risk.
4. Treat the strongest reasonable counterargument fairly.
5. State where the argument stops applying.
6. End with the practical position or workflow Morten currently holds.

Use conversational sentence rhythm, natural contractions, specific nouns, occasional parenthetical asides, and restrained humor. Correct spelling and grammar rather than copying mistakes from old posts. Preserve candidness, uncertainty, and rough edges in the thinking.

Never invent personal experience, client outcomes, implementation details, metrics, quotations, or opinions. Never merge separate published experiences into a new first-person story. Ask for missing first-person material when it is central to the approved angle.

### Avoid LLM voice

Do not use formulaic openings or filler such as:

- "In today's fast-paced world"
- "In the ever-evolving landscape"
- "Let's delve into"
- "game-changer"
- "unlock the power of"
- "seamlessly"
- "robust solution"
- "It is important to note"
- repeated "Whether you're..." constructions
- summaries that merely repeat the preceding section
- contrastive reversals such as "It is not X. It is Y."
- paired punchlines such as "X got cheap. Y did not."
- question fragments such as "The result?" or "The real problem?"
- empty transitions such as "That is where X comes in"

Do not use the X-but-Y pattern as a rhetorical cadence, even once. This includes adjacent sentences where the second sentence mainly negates, limits, or contrasts the first to create a punchline. Explain the relationship as part of a developed argument instead.

Introductions, body paragraphs, and transitional paragraphs must be developed across at least three complete sentences. Short paragraphs that present one compact claim or contrast are forbidden, even when they contain two sentences. The final paragraph of a section may contain one sentence when it genuinely concludes that section. Merge every other isolated or underdeveloped paragraph into the surrounding explanation rather than leaving it as standalone emphasis.

Do not use inflated metaphors, constant rhetorical questions, fake quotations, excessive bold text, or neat three-item lists when the material does not naturally have three parts. Vary sentence and paragraph length without becoming choppy.

### Technical accessibility

Write for technical people who are not necessarily developers:

- define a specialist term the first time it matters
- explain what a technical mechanism changes for the business, user, risk, cost, or delivery timeline
- include implementation detail only when it helps the reader evaluate or act
- use code only when the search intent genuinely requires it
- prefer a concrete example over a stack of abstractions
- never talk down to the reader

### Search and conversion quality

- Write to satisfy the approved primary keyword's intent. Do not let a broader topic replace it during drafting.
- Use the primary target keyword naturally in the title, first 150 words, and an appropriate heading when it reads well.
- Use supporting keywords only in sections that answer the corresponding reader need. Do not add sections or repeat phrases to reach a keyword quota.
- Keep the wording natural, but verify that every supplied keyword is accounted for according to the approved keyword decision.
- Link internally with descriptive anchor text.
- Cite external sources inline where a claim depends on them.
- Keep the CTA proportional. Connect MBV Labs only where the article naturally demonstrates relevant judgment or implementation experience.
- Prefer `/contact`, a relevant `/services`, `/work`, or `/projects` page over a generic homepage link.

## Draft Output

Return:

1. **Title**
2. **Slug**
3. **Primary target keyword**
4. **Excerpt** of one or two specific sentences
5. **Tags** with only useful topic labels
6. **Body** as clean Markdown, beginning with the introduction and no duplicate H1
7. **Source notes** listing factual claims and their URLs
8. **Review notes** with only unresolved factual or voice concerns

Before presenting it, check:

- the approved primary target keyword, search intent, and outline were followed
- every supplied keyword was handled according to the approved keyword decision
- every claim is sourced, clearly an opinion, or grounded in supplied experience
- the article adds something competitors did not
- a non-developer technical reader can follow it
- the voice resembles the strongest local samples
- headings are descriptive and scannable
- keywords are natural rather than repeated mechanically
- the CTA is relevant and restrained
- no prohibited filler or invented detail remains
- every X-but-Y reversal and paired punchline has been rewritten as a developed explanation
- every introduction, body paragraph, and transition has at least three complete sentences, with one-sentence paragraphs used only to conclude a section

Then stop and ask whether the user wants revisions or approves saving the draft through `/api/articles`.

## Approval Gate 2

Save only after explicit approval given after the complete draft is shown. Treat requests for revisions as withholding approval. Never infer save approval from outline approval, silence, or a general request to create an article.

Before sending, show a compact summary of the exact title, slug, excerpt, and tags that will be saved. If the user changes any field, update it and ask for save approval again.

## Save Draft

Send this shape to `POST /api/articles`:

```json
{
  "title": "Approved title",
  "slug": "approved-slug",
  "excerpt": "Approved excerpt",
  "body": "Approved Markdown without an H1",
  "tags": ["approved", "tags"]
}
```

The endpoint always creates a draft. Do not claim it is published. On success, report the returned article ID, title, slug, and draft status.

On a `409` slug conflict, do not silently alter the slug. Ask whether to use a proposed alternative. On validation or upstream errors, report the safe error details without exposing credentials and do not retry mutating requests blindly.
