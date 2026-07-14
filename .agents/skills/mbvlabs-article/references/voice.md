# Morten Vistisen Voice Guide

Use this as the complete tone-of-voice reference. It already captures the recurring patterns, strengths, and rough edges needed for article writing. No external voice research is required.

## The Voice in One Sentence

A working software engineer talks through how his opinion changed, grounds it in something he built or observed, pushes back on fashionable certainty, and ends with the practical position he currently holds.

## Core Traits

### 1. Experience before doctrine

Morten usually starts with something that happened:

- a previous article or opinion
- a project he built
- a contract or technical decision
- a failed attempt
- a change in his workflow

The experience earns the argument. Do not open with an abstract industry overview when a concrete starting point exists.

Typical movement:

```text
What I thought before
        |
What happened in practice
        |
Why the old view stopped holding
        |
What I do now
```

### 2. Strong opinions with visible uncertainty

The writing takes a position without pretending the position is permanent or universal. Morten regularly admits:

- his view changed
- the jury is still out
- a choice may turn out to be wrong
- another team could reasonably choose differently
- his own skill or context may affect the result

This is not weak hedging. State the argument clearly, then name its real boundary.

Good shape:

> I ended up choosing Pulumi because the team already knew Go. If we had a dedicated infrastructure team, I would probably have chosen Terraform.

Bad shape:

> Pulumi is the best infrastructure tool for modern teams.

### 3. Concrete technical detail tied to a consequence

The writing names real tools, code, failure modes, and systems. It then explains why the detail mattered:

- unfamiliar infrastructure code created team ownership risk
- generated indirection made production debugging harder
- a sleeping database broke onboarding
- a five-day deadline removed speculative features

Do not leave technical detail floating on its own. Connect it to understanding, ownership, delivery, reliability, cost, or business risk.

### 4. Conversational resistance to hype

Morten is skeptical of polished industry narratives, especially claims about AI productivity, best practices, and universal tooling choices. He often answers the obvious counterargument directly.

Useful techniques:

- name the fashionable claim
- test it against a real project
- ask what the claimed speed is actually producing
- distinguish a demo from software somebody must own
- acknowledge the strongest counterexample

Do not turn this into automatic contrarianism. The point is practical scrutiny, not taking the opposite side for attention.

### 5. Humor comes from irritation or specificity

The humor is dry, occasionally profane, and often parenthetical. It appears as a release valve in an otherwise serious explanation.

Examples of the pattern, not phrases to copy:

- a deliberately grand name for an ordinary industry actor
- a short aside about an absurd cost or claim
- self-deprecation after making a strong statement
- an unexpected cultural reference that fits the point

Use at most a few moments in a full article. Do not manufacture jokes, internet slang, or profanity to imitate personality. MBV Labs articles should usually use a slightly more restrained version of this trait than personal posts.

### 6. The prose shows the thinking process

Transitions often reveal the path rather than hiding it:

- `So, how did I arrive back where I started?`
- `But then something changed.`
- `Anyway, back to...`
- `The jury is still out on this one.`
- `I hear you, but...`
- `This is the tradeoff.`

These are evidence of a voice that reasons in public. Use the function, not the exact phrase. A draft should not sound like a fully packaged corporate conclusion existed before the experience.

## Sentence and Paragraph Rhythm

- Use contractions naturally.
- Prefer medium-length conversational sentences, then use a short sentence to land an important point.
- Let some sentences begin with `But`, `And`, or `So` when it improves the flow.
- Use parentheses for honest asides, not constant jokes.
- Introductions, body paragraphs, and transitions contain at least three complete sentences.
- Short paragraphs built from one compact claim or contrast are forbidden, even when they contain two sentences.
- A one-sentence paragraph is allowed only as the final paragraph that genuinely concludes a section.
- Lists appear when there is a real sequence, set of criteria, or summary. Do not force every argument into three bullets.
- Rhetorical questions are usually objections or decision questions, not engagement tricks.

The writing is conversational, not choppy. Do not imitate an LLM's habit of putting every sentence on its own line.

## Vocabulary and Register

Natural register:

- ordinary verbs such as `use`, `build`, `write`, `fix`, `choose`, and `ship`
- specific technical nouns
- `I think`, `I suspect`, `from my experience`, and `to be honest` when uncertainty is real
- `simple`, `practical`, `familiar`, `overkill`, `friction`, `tradeoff`, `ownership`, and `understanding`

Use these because the thought calls for them, not as a word-frequency costume.

Avoid replacing ordinary words with consultancy language. Morten writes `use`, not `leverage`; `help`, not `facilitate`; `hard to change`, not `creates transformation friction`.

## Grammar and Rough Edges

Published examples contain spelling mistakes, repeated words, and imperfect sentences. Those mistakes are not the voice.

Preserve:

- candidness
- direct transitions
- changing opinions
- uneven but natural paragraph rhythm
- occasional parenthetical asides
- willingness to say a tool or idea was frustrating, absurd, or wrong

Correct:

- spelling
- accidental repetition
- broken grammar
- unclear references
- factual errors

The target is Morten after a careful edit, not a sanitized corporate writer and not a reproduction of old typos.

Use British English spelling where variants differ. Do not use em dashes.

## Anti-LLM Cadence

A draft can avoid obvious buzzwords and still sound generated. The most common giveaway is artificial contrast written as a short punchline.

Treat these constructions as high-risk:

- `It is not X. It is Y.`
- `This is not about X, but about Y.`
- `The question is not X. The question is Y.`
- `X got cheap. Y did not.`
- `The first version is easy. The next five years are not.`
- `It does not just do X. It does Y.`
- `The real problem? X.`
- `The result? X.`
- `That is where X comes in.`
- `And that matters.`
- `Simple as that.`

The problem is not the word `but`. Natural writing uses contrast. The problem is repeated, symmetrical reversal that sounds engineered for a social-media quote.

Prefer one direct statement that explains the relationship:

```text
Generated:
The code got cheap. Maintenance did not.

Better:
AI can reduce the cost of the first implementation without reducing the work required to maintain it.
```

```text
Generated:
This is not a coding problem. It is an ownership problem.

Better:
The difficult part begins when nobody owns the generated application after it reaches production.
```

Also avoid:

- using one-sentence paragraphs for openings, body points, transitions, or artificial drama
- turning every section ending into a slogan
- repeated sentence fragments that restate the previous sentence
- perfectly balanced three-part lists when the material does not naturally have three parts
- asking a rhetorical question and immediately answering it throughout the article
- repeated `Here is the thing`, `The catch`, `The truth`, or `The reality` transitions
- generic aphorisms that could appear in an article on any subject
- making every paragraph the same length
- ending several sections with a neat moral

Introductions, body paragraphs, and transitions must contain at least three complete sentences. Short paragraphs built from one compact claim or contrast are forbidden, even when they contain two sentences. A section's final paragraph may contain one sentence when it genuinely concludes the section. Merge every other isolated or underdeveloped paragraph into the surrounding explanation.

During the final voice pass, search specifically for contrastive reversal and punchline fragments. Rewrite every compact X-but-Y reversal, including adjacent sentences where the second mainly negates or limits the first for emphasis. Natural contrast may remain only when it is developed as part of an explanatory paragraph rather than presented as a punchline.

## Article-Specific Voice Shapes

### Opinion or changing-view article

Use this when the input is an argument, reaction, or evolving belief.

1. Start with the previous belief or recent trigger.
2. State the current tension early.
3. Explain what happened in practice.
4. Show the cost or failure mode with a real example.
5. Address the strongest reasonable counterargument.
6. Explain the narrower role where the disputed tool or practice still helps.
7. End with the workflow or position Morten currently uses.

The conclusion should be provisional but useful: `This is what I do now`, not `The industry must do this forever`.

### Comparison or decision article

Use this when choosing between tools, approaches, or build options.

1. Open with the actual decision and team context.
2. Explain why the decision mattered beyond technical preference.
3. Introduce each option only as deeply as the decision requires.
4. Describe a proof of concept, test, or concrete comparison.
5. Name the criteria that determined the choice.
6. State the choice.
7. Give the conditions under which the other option would win.
8. Admit what remains unknown and when the decision should be revisited.

The point is not to crown a universal winner. It is to help a reader make the same class of decision in their own context.

### Build log or project retrospective

Use this when writing from a challenge, launch, or implementation story.

1. Explain the personal reason for attempting it.
2. Define the constraint and initial plan.
3. Move through two or three concrete failures or surprises.
4. Show enough technical detail to make each problem real.
5. Report the honest result, including what did not ship.
6. Extract one or two lessons that apply beyond the project.

Do not clean up the story until it sounds inevitable. The failed assumptions are usually where the value lives.

### Technical guide

Use this when the reader needs to implement something.

1. State who the guide is for and what familiarity it assumes.
2. Explain the practical system being built, not a toy abstraction when a real example is available.
3. List only required tools and prerequisites.
4. Build in a sequence the reader can run and verify.
5. Explain important code immediately after it appears.
6. Tie structural choices to testability, readability, adaptability, or ownership.
7. Mark shortcuts and outdated details honestly.
8. End with a working result, limitations, and relevant resources.

The code carries much of the technical weight. Prose should explain decisions and consequences rather than narrating every line.

## MBV Labs Adaptation

The personal site voice and the MBV Labs business goal must coexist.

Keep from the personal writing:

- first-person experience
- practical skepticism
- honest tradeoffs
- concrete technical nouns
- restrained humor
- a conclusion grounded in current practice

Adjust for MBV Labs:

- make the buyer's decision or business pressure explicit
- explain specialist terms when founders or product leaders may not know them
- remove profanity unless the user explicitly wants it
- keep cultural references accessible
- connect the technical point to delivery, risk, cost, ownership, or maintainability
- use one light conversion path near the end

Do not insert MBV Labs into every section. The article should first earn trust by being useful.

## Applying This Guide

Before outlining:

1. Choose the closest article shape from this guide.
2. State which reasoning moves fit the topic, such as changing view, project evidence, fair counterargument, or bounded recommendation.
3. Identify the personal story, opinion, result, or client context needed to support those moves.
4. Use only facts supplied in the request, approved local context, or diary entries fetched for the current task.
5. Ask Morten for missing first-person material before drafting.

Add a `Voice plan` to the research brief. It should name the planned opening, source of first-person authority, counterargument, uncertainty boundary, humor level, and conclusion style. It should not name or retrieve old source articles.

Never combine details from several sources into a new first-person event. Never write `I found`, `I built`, `I prefer`, or `I tell clients` unless material approved for the current task supports it.

## Final Voice Check

Read the draft once for facts, then once only for voice.

Ask:

- Does this begin from a real situation rather than an SEO preamble?
- Is the main opinion clear without sounding universal?
- Can the reader see why Morten holds this view?
- Is there at least one concrete example doing real argumentative work?
- Is the strongest counterargument treated fairly?
- Does technical detail lead to a practical consequence?
- Are uncertainty and limitations specific rather than generic hedging?
- Does the conclusion state a useful current position?
- Did editing remove mistakes without removing personality?
- Have all compact X-but-Y reversals and punchline fragments been rewritten as developed explanations?
- Does every introduction, body paragraph, and transition contain at least three complete sentences?
- Is every one-sentence paragraph the final paragraph of its section and a genuine conclusion?
- Would this still sound plausible if the MBV Labs name were removed?
