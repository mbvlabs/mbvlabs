# MBV Labs Search Article Structure

This guide separates the useful search architecture seen in Outrank's writing examples from the parts that would make an MBV Labs article feel automated or overproduced.

Examples examined include:

- `How to Do a Website Audit That Actually Boosts Your SEO`
- `Mastering the Perfect Blog Post Format`
- `10 Best Semrush Alternatives Compared in 2026`
- `How Many Referring Domains Do You Need to Rank?`
- `How AI Reads Your Website`

Outrank's examples are evidence for search coverage and page architecture. They are not voice references.

## What to Borrow from Outrank

### Answer the query early

Strong search articles do not make readers wait through a broad history lesson before receiving the basic answer.

The introduction should normally establish:

1. the reader's problem or decision
2. the direct answer, thesis, or useful rule of thumb
3. why the obvious answer is incomplete
4. what framework, steps, or comparison follows

For a narrow question, give the answer in the first 150 words. For an opinion article, make the tension and current position clear in the opening.

### Build headings from reader jobs

Outrank's useful headings correspond to things a reader needs to understand or do:

- define the concept
- understand why it matters
- compare options
- calculate a target
- follow a process
- recognise failure cases
- choose a next step

Each H2 should resolve one distinct reader job. H3s break down that job only when it contains several meaningful parts.

### Cover the decision, not merely the definition

The stronger examples move beyond `what is X` into:

- when X matters
- how to evaluate X
- how much X is enough
- when X will not solve the problem
- what X costs
- which option fits which situation

This is particularly valuable for MBV Labs because decision support demonstrates technical leadership better than a dictionary-style explainer.

### Use information-dense formats

Choose the format that communicates the material fastest:

- comparison table for options or criteria
- numbered steps for a real sequence
- checklist for verification
- code for implementation
- decision matrix for tradeoffs
- example scenarios for applying a framework
- simple ASCII diagram for architecture or flow

Do not turn normal prose into a table merely for visual variety.

### Close unresolved search intent

A short FAQ can capture genuine questions from People Also Ask or repeated competitor coverage. Include it only when the main article would become awkward if every side question were forced into the core flow.

Answers should be concise. Do not repeat whole sections in FAQ form.

### Create an internal path

Useful Outrank articles connect related pages throughout the explanation. MBV Labs should do the same more selectively:

- link a service when the reader reaches a decision MBV Labs can help execute
- link a work example when it proves the claim
- link a project when it gives the reader something concrete to inspect
- link another article when it covers a necessary tangent in depth
- link `/contact` once the article has established a relevant problem

Anchor text should describe the destination. Do not stuff exact-match keywords into every link.

## What Not to Borrow

Several Outrank examples run between roughly 3,000 and 4,300 words. Length is not the useful pattern. Do not copy:

- a fixed long-form word target
- an H3 under every H2
- repeated explanations of the same benefit
- generic analogies in every section
- unsupported statistics used as decoration
- false first-person authority such as `I see this all the time`
- headings added only to capture a keyword variation
- an FAQ that restates the body
- multiple promotional sections before the conclusion
- a CTA longer than the useful content section before it
- generic SEO phrases such as `moves the needle`, `secret sauce`, or `game-changing`
- forced historical background when the reader came for a practical answer

MBV Labs should be more selective, more technical, and more candid.

## The Core Article Spine

Every article needs a spine, not a rigid template.

```text
Search query or reader problem
              |
Direct answer or Morten's current thesis
              |
Context needed to understand the decision
              |
Evidence, implementation, comparison, or framework
              |
Limits, tradeoffs, and failure cases
              |
Practical next step
```

The approved outline should make this progression visible. If two adjacent sections do the same job, merge them.

## Adaptive Structures

Select one primary structure based on search intent and voice evidence. Add sections only when research shows they are needed.

### 1. Practical guide

Use for `how to`, implementation, audit, migration, and workflow queries.

```markdown
Introduction
- Name the real situation
- Give the short answer
- State prerequisites and expected result

## Before you start
- Assumptions
- Required access or tools
- Important constraint

## Step or phase one
- Action
- Why it matters
- Verification

## Step or phase two
- Action
- Why it matters
- Verification

## The tradeoff or failure mode most guides miss
- Concrete problem
- How to diagnose it
- When to choose another approach

## What the finished version should look like
- Result
- Checklist
- Known limitation

## A practical next step
- Conclusion
- Restrained CTA
```

For code-heavy guides, organise around working milestones rather than arbitrary topic headings. Every major step should leave the reader with something they can verify.

### 2. Decision framework

Use for `should I`, `when to build`, architecture, buy-versus-build, and operational decisions.

```markdown
Introduction
- State the decision
- Give Morten's short answer
- Explain why context changes the answer

## What is actually being decided
- Business pressure
- Technical constraint
- Ownership horizon

## The criteria that matter
- Criterion with consequence
- Criterion with consequence
- Criterion with consequence

## Where option A works
- Best-fit conditions
- Cost or risk

## Where option B works
- Best-fit conditions
- Cost or risk

## The failure mode people underestimate
- Real example or scenario

## A simple way to choose
- Decision checklist or matrix

## What I would do in this situation
- Bounded recommendation
- Natural MBV Labs next step
```

This structure should help the reader make a decision, not quietly force every reader toward custom software or consulting.

### 3. Comparison

Use for `X vs Y`, alternatives, tools, frameworks, and platforms.

```markdown
Introduction
- Actual context for the comparison
- Recommended option by situation
- Main tradeoff

## The short version
- Compact comparison table

## What matters for this decision
- Evaluation criteria tied to team and business context

## Option A in practice
- Relevant strengths
- Relevant limitations
- Evidence or proof of concept

## Option B in practice
- Relevant strengths
- Relevant limitations
- Evidence or proof of concept

## Where the difference showed up
- Real test, implementation, or operating scenario

## Why I chose one
- Choice for the stated context

## When I would choose the other
- Fair reversal conditions

## What remains unknown
- Review point or unresolved risk
```

Avoid ranking ten products when research and firsthand knowledge only support comparing two or three.

### 4. Opinion backed by experience

Use for contrarian arguments, changing views, and industry commentary.

```markdown
Introduction
- Previous belief or trigger
- Current position

## What changed
- Experience or evidence

## Why the common argument is incomplete
- Industry claim
- Practical problem

## Where it broke in a real project
- Concrete story
- Technical consequence
- Business consequence

## The strongest counterargument
- Treat it fairly
- State the boundary of Morten's claim

## Where the idea still helps
- Narrow useful role

## The workflow I use now
- Current practical position

Conclusion
- One clear takeaway
- Light conversion path if relevant
```

Search optimisation should not flatten the narrative. Supporting queries can shape headings where they fit, but the experience remains the spine.

### 5. Build log or retrospective

Use for project launches, experiments, migrations, and time-boxed builds.

```markdown
Introduction
- Why attempt this
- Constraint
- Honest result upfront

## The initial plan
- Scope
- Architecture or approach

## First surprise or failure
- What happened
- Diagnosis
- Change

## Second surprise or failure
- What happened
- Diagnosis
- Change

## What shipped
- Completed work
- Missing work

## What I would do differently
- Specific lesson

## The broader takeaway
- Relevance to founders or teams
- Natural next step
```

Do not remove the failed paths to make the project look cleaner. They are often the most useful and credible parts.

### 6. Focused answer article

Use for a narrow long-tail question that does not need a large guide.

```markdown
Introduction
- Direct answer in the first paragraph
- One qualification

## Why the answer depends on context
- Two or three deciding factors

## A practical way to calculate or choose
- Small framework or example

## When the rule fails
- Boundary or failure case

## What to do next
- Concise conclusion
```

This may be 700 words or 1,500 words. Do not expand it to compete with unrelated 4,000-word pages.

## Introduction Design

Write the introduction after the body outline is stable.

A useful MBV Labs introduction usually contains four moves:

```text
Concrete situation
      +
Direct answer or tension
      +
Why common advice misses part of the problem
      +
What this article will resolve
```

Prefer a real situation over a formulaic hook. Do not use Problem-Agitate-Solve mechanically. The reader does not need their pain dramatized if they already searched for a technical answer.

Avoid:

- dictionary definitions as the first sentence
- `In today's...`
- cinematic scene setting
- unsupported claims about what `most businesses` do
- telling readers the article is `comprehensive` or `ultimate`
- explaining why the topic is important for several paragraphs before answering

## Section Design

Introductions, body paragraphs, and transitions must contain at least three complete sentences. Short paragraphs built from one compact claim or contrast are forbidden, even when they contain two sentences. A one-sentence paragraph is allowed only as the final paragraph that genuinely concludes a section.

Every section should contain at least two of these:

- a clear claim
- an explanation of the mechanism
- a concrete example
- evidence or a source
- a tradeoff
- an action or decision

A section that contains only generic explanation is probably padding.

For technical sections, use this small pattern:

```text
What it is
   |
How it works
   |
Why it matters here
   |
What can go wrong
```

For non-developer technical readers, define only the term needed to understand the consequence. Do not turn the article into a glossary.

## Evidence Architecture

During outlining, tag planned material as one of:

- `[experience]` supplied by Morten or published in his writing
- `[primary source]` official docs, standards, or original research
- `[live SERP]` observable result-page pattern
- `[example]` clearly presented hypothetical scenario
- `[opinion]` Morten's bounded judgment

Do not draft an unsupported factual section and plan to find a citation later. If evidence is missing at outline time, mark it as a question for Morten or remove the claim.

Competitor articles can reveal a topic to verify. They are not sufficient evidence for a statistic merely because several repeat it.

## FAQ Rules

Include an FAQ only if at least two useful questions remain after the main outline.

A question belongs when it:

- appears in People Also Ask or related searches
- affects a reader's decision
- can be answered accurately in a short paragraph
- does not deserve a full core section

Use the exact search wording only when it sounds natural. Three useful questions are better than eight thin ones.

## Conversion Structure

The article should produce clients by demonstrating judgment, not by interrupting the reader repeatedly.

Choose one primary conversion bridge:

- `Need help deciding?` for decision articles
- `Need this implemented?` for guides
- a relevant work example for proof-led articles
- a project page when the tool is inspectable
- `/contact` when the reader is likely facing the described problem now

The final CTA should normally be one short paragraph. It should name the situation in which contacting MBV Labs makes sense. Do not claim the article has created urgency when it has not.

## Outline Quality Check

Before requesting approval, verify:

- Can the primary query be answered from the proposed title and first two sections?
- Does every H2 perform a different reader job?
- Is the article type clear?
- Is Morten's experience or judgment doing work that competitors cannot copy?
- Are technical details tied to practical consequences?
- Is evidence assigned before drafting?
- Is there a fair limitations or counterargument section?
- Are internal links relevant rather than quota-driven?
- Is the FAQ optional rather than automatic?
- Is the conversion path proportional to the article?
- Could any section be deleted without harming the reader? If yes, delete it now.
