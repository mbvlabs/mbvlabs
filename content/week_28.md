# Content for Week 28

**Week:** July 6 to July 12, 2026  
**Business:** MBV Labs

## Voice notes

These drafts are using this week's diary entries as the voice reference. The tone is personal, technical and still figuring things out in public. Keep the rough edges when publishing. It should feel like a person arriving at an opinion, not a brand presenting a finished doctrine.

The main thought running through the week is that AI is making code much easier to produce, which leaves me wanting to spend more time deciding what should be built in the first place. The other threads around UI work, internal tools, DeployCrate and open source all come back to that in one way or another.

---

## X / Twitter posts

### Post 1

I’m spending less time writing code and more time doodling on paper before I open the editor.

Once I understand the problem, getting an agent to follow the recipe is almost trivial.

Getting to the right recipe is still the work.

### Post 2

AI can write more code than I could ever review.

So producing code is a pretty useless metric now.

I’m much more interested in how much I can remove before anything gets built. Usually the product gets better when I do.

Cut that, cut that, cut that.

### Post 3

My current split with AI coding is roughly:

I think through the problem, write the parts I care about and let the clanker handle the predictable UI.

Some CRUD table with filters and a modal? Go nuts.

Core domain logic that I’ll live with for years? I’m staying close to that.

### Post 4

I spent years writing React, grew to hate it, moved back towards server-rendered HTML and somehow AI has brought me all the way back around to SPAs.

Vue + Inertia has been a bit of a revelation.

The backend keeps the logic. The clanker gets a UI stack it actually understands. I get to work on the bits I enjoy.

### Post 5

One thing I’ve been missing with AI coding: the model knowing your stack matters a lot.

I kept fighting it on templ + Datastar. Give it Vue, Inertia and shadcn-style components and it can hammer out a decent admin UI in an evening.

That is starting to influence how I choose tools.

### Post 6

I’m not comfortable with the “I haven’t written a line of code all year” thing.

Coding is a skill. Stop using it and eventually you lose some of it, including the part that tells you when generated code is rubbish.

I want the leverage. I also want to keep my hands in the work.

### Post 7

Replacing SaaS with an AI-written internal tool looks great when you only count the first afternoon.

A year later you have 25 tiny apps, 25 sets of dependencies and a pile of software nobody really owns.

The code got cheap. Maintenance didn’t.

### Post 8

I’m starting to think deployment is the less interesting part of DeployCrate.

Most people can get an app online somehow.

Knowing what is happening once it is online is still a mess. The data exists, but turning logs and metrics into an answer a builder can act on feels much less solved.

### Post 9

I’ve decided to open up DeployCrate, which was actually my first instinct.

Developers want to inspect things. They want to run them and decide for themselves before convincing a boss to pay for a hosted version.

I let advice from outside that market pull me away from something I already understood about it.

### Post 10

Founders usually aren’t looking for someone to close a neat stack of tickets.

They have a fuzzy product or technical problem and need someone senior to make sense of it, cut it down and then actually ship the thing.

That combination is what I’m building MBV Labs around.

### Thread

**1/6** My view on AI coding keeps moving around. I went from completely sold, to quite sceptical, and now I’m landing somewhere in the middle.

**2/6** The models are extremely good at familiar, repetitive work. I gave one a Vue + Inertia admin UI with known components and got most of an evening’s work back.

**3/6** They are much less useful when the direction is fuzzy. You get a lot of code very quickly, then spend your time discovering that it solves the wrong version of the problem.

**4/6** So I’m changing where I spend my time. More thinking, sketching, cutting scope and working out the important boundaries. Less hand-writing every table and modal.

**5/6** I’m still writing code because I want to keep the skill that lets me review what comes back. Fully removing myself from implementation feels risky, especially when the current pricing and subscriptions could change at any point.

**6/6** This is also how I’m approaching MBV Labs work. Stay close enough to the business to choose the right thing, and close enough to the code to ship it properly.

---

## LinkedIn posts

### Post 1: I’m changing how I spend my coding time (QUEUED)

My thinking on AI coding keeps changing.

I was an early adopter, then became quite sceptical, and now I’m settling into a middle ground that feels more useful than either extreme.

I’m spending more time away from the editor. Usually with paper, trying to understand the problem and work out which parts are essential. Once I have that, an agent can follow the recipe very quickly.

The difference has become obvious while I’ve been building an admin UI for my personal site. With Vue, Inertia and components the model already knows, it produced most of the interface in an evening. It looks better than the UI I replaced and I didn’t have to spend hours on a part of the work I don’t particularly enjoy.

That is a good use of the tool for me. I still own the backend, the data model and the decisions that will be painful to change later. The model handles a lot of the predictable interface work, and I review what it gives me.

I’m still seeing people measure AI progress by how little code they write. That makes me uneasy. If I stop coding completely, I’m also going to lose some of the skill I need to spot bad generated code. There is also a fairly large dependency on subscriptions staying cheap enough for this whole setup to make sense.

So I’m keeping my hands in the code. I’m just becoming much more selective about which code deserves my time.

For the work I do through MBV Labs, that split feels right. A founder needs somebody thinking through the product and technical tradeoffs with them, then staying around to get the software shipped. Generating a bigger pile of code was never the useful part.

### Post 2: The internal tool hangover (USED)

I was talking with a friend about people replacing SaaS products with small tools they have built using AI.

I can see why it is happening. A subscription looks expensive next to an app an agent can produce in an afternoon, especially when you are only using a fraction of the SaaS product anyway.

I’m not sure the pattern is going to age particularly well.

One tool is fine. Then every team does it and suddenly the company has 20 or 30 little applications. They all have dependencies. They all need authentication, permissions, backups and monitoring. Somebody needs to update them when an API changes or when the person who prompted the first version leaves.

By then you have built an internal platform by accident.

There are obviously good reasons to own software. If a workflow is particular to how the company operates, and getting it right gives the business an advantage, custom software can be the sensible option. I build these kinds of tools for clients.

But I’d want us to be honest about the whole decision. Who is looking after it next year? What happens when it breaks on a Friday? Is the workflow valuable enough that we want to own all of that?

AI is bringing the first version close to zero in some cases. It is not carrying the pager for you afterwards.

### Post 3: AI has somehow brought me back to SPAs (USED)

I used to write a lot of React. It was most of what I did in my first two jobs and I genuinely loved it for a while.

Then everything around it seemed to get more complicated. For my own projects I moved back towards Go, templ, HTMX and Datastar. Those tools fit how I like to build systems, and I enjoyed working with them much more.

The UI was always the weak part. I’m much more comfortable on the backend and my designs tend to be derivative. They were fine, but never particularly good.

I also couldn’t get AI coding tools to do much useful frontend work in that stack. The models didn’t know it well enough and would make things up. I spent more time correcting them than I saved.

While finishing Andurel, somebody suggested adding Inertia support. I rejected it at first because it felt like it went against the whole idea of the framework. Eventually I tried it with Vue.

I’ve been missing out.

I can keep routing and business logic in Go, keep the MVC shape I like, and give the model a frontend stack it knows extremely well. It can now get surprisingly close on the first attempt, especially when I point it at a known component library.

This is changing how I’m thinking about stack selection. I still care about simplicity, maintainability and whether the technology fits the system. Now I’m also considering how well an agent can work with it, particularly for routine UI that is slowing down the rest of the product.

I’m even rewriting my personal blog with React and Inertia next, which is not a sentence I expected to write a few months ago.

### Post 4: I’m opening DeployCrate (QUEUED)

I’ve decided to make DeployCrate open source or source available.

That was my original plan. I moved away from it after a successful entrepreneur convinced me that keeping it closed was the better route.

The advice made sense from his point of view, but he isn’t selling to developers. I am, and I should have trusted what I knew about that audience.

Developers want to look around before they commit to infrastructure. They want to inspect the code, run it themselves and understand how it behaves. If it works for them, they can take it into a company and make the case for paying for a hosted version.

Keeping everything closed meant I was asking them for trust before giving them much of a way to build it.

There is also a product question I’m still working through. Deployment is largely solved now. There are many good platforms and plenty of ways to have an agent spin something up. The part that keeps bothering me is what happens after deployment.

Telemetry is still hard to make sense of. Logs and metrics are everywhere, but a builder often finds out they lacked the right insight after something has already broken. I’m wondering if DeployCrate should simplify the deployment side, lean on buildpacks and put much more energy into understandable telemetry, backups and the things that help someone operate an application.

I don’t have the positioning neatly wrapped up yet. Opening the code feels like the right next move though. It gives developers a way to try the product on their own terms, and it gives me a much better feedback loop while I narrow the rest.

Sometimes you spend a while going around the houses only to end up back at your first instinct, hopefully with enough new information to make it work this time.

---

## Blog post drafts

### Draft 1: I’m Writing Less Code, but I’m Not Giving Up Coding

**Suggested slug:** `/blog/writing-less-code-with-ai`  
**For:** Founders and small technical teams

My thinking on AI coding has moved around quite a lot.

I started out completely sold on it. Then I became more sceptical, partly because I could see the weaknesses more clearly and partly because the culture around it became a bit breathless. Now I’m trying to find a middle ground where the tools are useful without pretending they can take responsibility for the entire job.

The change I’m making is fairly simple. I’m spending more time thinking and less time typing.

Before I open the editor, I’m doodling on paper and trying to understand the problem space. What is the smallest version that does something useful? Which decisions will be hard to undo? Which code expresses the actual business and which code is just another table, form or integration that follows a known pattern?

Once that path is clear, getting an agent to follow it is becoming almost trivial.

#### The model does much better when the road already exists

I saw this while building the admin area for my personal site. Frontend work has always slowed me down. I can do it, and I spent years writing React professionally, but I’ve always gravitated towards backend systems. My own interfaces tended to be okay and rarely much better than that.

I had also struggled to get models to work well with my preferred stack of Go templates and less common frontend tools. They would hallucinate APIs or produce code that looked plausible until I tried to use it.

Then I added Inertia support to Andurel and tried Vue with a component library the model knew. In one evening it produced nearly the whole admin UI. The result was better than the interface I was replacing.

That experience helped me find a sensible boundary. I’m happy to let the model build predictable UI when I have already chosen the shape of the system. I’m keeping close control of the backend, domain model, permissions and the odd parts where the product actually differs from everything else.

Some list with CRUD actions does not need to be my artistic expression. The clanker can have it.

#### Cheap code makes cutting more important

There is a trap here. When a feature becomes cheap to produce, every feature starts to look worth producing.

Previously, the cost of design and engineering forced a conversation about scope. Now an agent can build the idea before the team has properly decided if it belongs in the product. You can end up with a larger application very quickly and discover later that most of the complexity came from things nobody really needed.

I’m finding myself borrowing a phrase from the Always Sunny podcast whenever I look at a product plan: cut that, cut that, cut that.

More generated code is not helping if it pushes the product in the wrong direction. I’d rather spend an hour removing a workflow than ten minutes prompting it into existence and several years carrying it around.

#### I still want to know how to code

I keep seeing people celebrate that they have not written a line of code all year. I understand the excitement, but I’m not seeing that as a goal for myself.

Coding is a skill like anything else. If I stop using it, I will lose some of it. That also weakens my ability to review what the model produces, recognise a poor abstraction or notice when a permission check is sitting in the wrong place.

There is a practical dependency too. This workflow currently makes sense because good models are available through subscriptions at a manageable price. If those subscriptions disappear or the economics change, I don’t want to discover that I have outsourced the whole craft to a tool I can no longer use in the same way.

So I’m still writing code. I’m saving that attention for the parts where my experience is useful and letting the model clear away more of the routine work.

#### This is where technical leadership is moving

For a founder, the difficult part of a software project is rarely getting somebody to type every line. It is working out what should exist, turning a fuzzy business need into a sensible system and making tradeoffs that the company can live with.

AI is making implementation faster, which is great. It is also making it much easier to build the wrong thing with impressive speed.

The work I’m doing through MBV Labs sits in that gap. I can help a founder understand the problem, choose a practical route and then stay close enough to the implementation to ship it. I’m using AI in that work, of course, but the useful outcome is still software that solves the right problem and can be maintained once I’m gone.

If you’re sitting on an important software project and the shape of it is still unclear, you can get in touch at [mbvlabs.com](https://mbvlabs.com).

---

### Draft 2: The First Version Is Cheap. Who Owns the Next Five Years? (USED)

**Suggested slug:** `/blog/ai-built-internal-tools`  
**For:** Founders, operations teams and engineering leads

I’ve been thinking about what happens when every company starts replacing SaaS products with small internal tools written by AI.

The appeal is obvious. You are paying for a large product and using perhaps ten percent of it. An agent can build the bit you need in an afternoon. The spreadsheet looks quite convincing at that point.

I suspect a lot of companies are going to wake up with an internal software estate they never meant to create.

One tool is easy to ignore. Then somebody builds another for finance, another for customer support and another to move data between the first two. A year later there are 25 of them. Each one has dependencies, credentials, user access and a database that matters to somebody.

The company has built a platform by accident, and nobody is quite sure who owns it.

#### The demo is a very small part of the cost

AI is genuinely reducing the time it takes to get a first version working. That does not tell us much about the next few years.

Somebody still needs to update dependencies, respond when an integration changes, keep the authentication working and know how to restore the data. The workflow itself will change because the business changes. The original prompt will not be around to explain why a strange decision was made six months ago.

This is where the buy versus build calculation gets distorted. We compare the monthly SaaS bill with the afternoon it took to generate a replacement. The maintenance work gets treated as if it will somehow take care of itself.

It rarely does. Usually it lands on the existing engineering team in small interruptions that are difficult to measure and very easy to underestimate.

#### I still think companies should build internal tools

Custom internal software can be extremely valuable. I build it for clients through MBV Labs, so I’m hardly going to argue otherwise.

The best cases tend to be workflows that are genuinely specific to the business. Perhaps the software connects systems in a way no off-the-shelf product supports, removes hours of manual work or captures something important about how the company operates. Owning that can be worth the maintenance because the software is doing more than avoiding a subscription.

Before building, I’d want to sit with the team and talk through a few fairly ordinary questions.

Who is looking after this next year? How often is the workflow changing? What happens if the tool is unavailable for a day? Can we get our data out if we decide to return to a vendor? Does owning this process actually help the business win?

Those questions usually tell us whether to buy a product, automate around one or build something ourselves. If nobody can name an owner, I’m already leaning towards waiting.

#### Build the things worth owning

AI is giving small teams the ability to make software that would previously have sat forever in a backlog. I’m excited about that. There are awkward, valuable workflows all over companies that finally have a chance of being fixed.

I just don’t want us to confuse the ability to generate an application with a reason to own one.

If your team is looking at a manual process or an expensive SaaS product and wondering whether custom software makes sense, I can help work through it. Sometimes the answer will be a focused internal tool. Sometimes it will be a small automation around what you already have. Both are better than inheriting a platform by accident.

You can tell me about the workflow at [mbvlabs.com](https://mbvlabs.com).

---

### Draft 3: Why I’m Opening DeployCrate

**Suggested slug:** `/blog/opening-deploycrate`  
**For:** Developers and people building software products

I’ve decided to make DeployCrate open source or source available.

This was actually the original plan. I moved away from it after talking to a successful entrepreneur who convinced me a closed product would be easier to turn into a business.

His advice was sensible for the world he knows. He is not a developer though, and he does not sell infrastructure to developers. I let general business advice override something I already understood about the particular audience I was trying to reach.

Developers like to inspect things.

They want to run the software themselves, look through the code and understand what it is doing before they put an application on it. If the experience is good, they are also in a much better position to convince a company to pay for a managed version.

With a closed product I was asking people to trust DeployCrate before they had much of a way to decide whether it deserved that trust.

#### I may have been focusing on the wrong part

There is another question sitting underneath the distribution model. I’m not sure deployment is the strongest reason for DeployCrate to exist.

Deployment is largely solved. There are many platforms doing it well, and an increasingly capable group of agents can produce the configuration for people who want to manage it themselves.

Operating the application afterwards is still awkward.

You can collect enormous amounts of telemetry, but turning it into an understandable answer is another thing. Builders often discover that they lacked the right logs or metrics after the problem has happened. Products like Datadog and Sentry are powerful, and can become expensive or overwhelming for a smaller company that wants to know why its application is misbehaving.

I’m wondering if DeployCrate should become much simpler around deployment. Use buildpacks, make backups straightforward and spend more energy translating telemetry into something a builder or less technical product owner can act on.

That direction needs more conversations and more evidence. I’m writing it down here because it is the part of the product I’m finding interesting right now.

#### Going around the houses was still useful

I could look at the detour and think I wasted time. I don’t.

DeployCrate pushed Andurel much further and exposed plenty of things I can remove from the next version. I’ve learned more about the product, the audience and where I’m willing to spend my time. The open version should end up simpler because of that work.

The next step is to choose the licensing model, open the code and let people use it in the way developers naturally evaluate software. From there I can see whether the telemetry idea holds up outside my own head.

I’m still lost on parts of the positioning, which is probably more honest than pretending a pivot arrives fully formed. I have a direction now and a better feedback loop. That is enough to keep moving.

I’ll be sharing the code and the decisions through [MBV Labs](https://mbvlabs.com) as it develops.

---

## YouTube video ideas

### Video 1: How I’m Actually Using AI to Code

Open on the paper sketch, before showing an editor. Talk through how the process is shifting towards more problem exploration and less hand-written routine code. Then build one admin screen with an agent, review what it produced and point out the decisions that had to happen before the prompt.

Possible titles: **How I’m Actually Using AI to Code in 2026**, **I Let AI Write the UI**, or **My AI Coding Workflow After the Hype**.

This can naturally connect to MBV Labs by showing what senior implementation looks like when AI is part of the toolkit. Keep it as a real build rather than a general talk about productivity.

### Video 2: AI Has Brought Me Back to SPAs

Tell the story from years of React, through Go templates and Datastar, to trying Vue with Inertia in Andurel. Build the same small interface in both styles and talk honestly about where the model struggled, where it did well and which stack you would choose for different products.

Possible titles: **Why AI Made Me Try SPAs Again**, **Go + Inertia + Vue Is Surprisingly Good**, or **I Hated React, Then AI Changed the Tradeoff**.

The interesting part is the full-circle story, so leave room for the personal history. A perfectly polished stack comparison would lose the thing that makes this one yours.

### Video 3: The SaaS Tool You Replace With AI Is Yours Forever

Start with a small SaaS bill and generate a rough replacement live. Once it works, zoom out and add everything missing from the demo: auth, permissions, backups, monitoring, updates and an owner. Use three real workflow examples to show when you would buy, automate around a product or build something custom.

Possible titles: **Before You Replace SaaS With AI**, **The Hidden Cost of AI-Built Internal Tools**, or **Your AI Tool Took an Afternoon. Now Maintain It.**

This is a strong MBV Labs topic because it shows business judgment alongside implementation. The conclusion can be nuanced. Some software is absolutely worth owning.

### Video 4: Why I’m Opening DeployCrate

Make this a direct build-in-public update. Explain the advice that moved the product away from open source, what developer behaviour has shown since, and why you’re returning to the original idea. Then walk through the unresolved part: deployment may be less valuable than making telemetry understandable.

Possible titles: **Why I’m Opening DeployCrate**, **I Took the Wrong Advice for My Developer Product**, or **DeployCrate Is Changing**.

Publish this once the license is decided so the video can end with something concrete people can inspect or follow.

### Video 5: Building an Admin UI With Go, Inertia, Vue and an Agent

Build a real resource screen with listing, creation, editing, validation and permissions. Keep the business logic in Go and show exactly what context you give the agent before it touches the frontend. Spend time reviewing the result, especially the places where visually correct code can still be wrong.

Possible titles: **Let AI Build the Frontend While Go Owns the Logic**, **Building an AI-Friendly Admin UI With Go**, or **Can an Agent One-Shot This Vue Admin Screen?**

This works as a longer tutorial for the MBV Labs channel and as public proof for Andurel. Use a real project or a believable slice of one so the tradeoffs are visible.

---

## Short-form drafts

### Short 1: I still write code

**Hook on screen:** I don’t want AI to stop me writing code.

I keep hearing people say they haven’t written a line of code all year.

I’m not seeing that as a goal.

Coding is a skill, and I need that skill to tell when the generated code is bad. I want AI clearing away the repetitive work while I stay close to the decisions that matter.

The leverage is useful. So is knowing what I’m looking at.

### Short 2: Your afternoon tool

**Hook on screen:** That internal tool took an afternoon. Now who owns it?

An agent can replace part of a SaaS product surprisingly quickly.

Then the tool needs auth, backups, monitoring, security updates and somebody who remembers why it works the way it does.

Do that 20 times and you’ve built a company platform by accident.

The first version is cheap. Make sure the software is worth owning afterwards.

### Short 3: Cut that

**Hook on screen:** My favourite product prompt right now.

AI can produce more features than anybody needs.

So whenever I’m looking at a product plan, I’m hearing the same phrase: cut that, cut that, cut that.

Remove the extra workflow. Remove the clever abstraction. Remove the screen that only exists because it was easy to generate.

Cheap code is making scope more important.

### Short 4: Give the model a stack it knows

**Hook on screen:** I thought the model was bad at frontend.

I kept asking agents to write UI with Go templates and less common tools. They kept making things up.

Then I tried Vue, Inertia and a component library they knew. In one evening I had most of an admin UI.

The tools your model knows are starting to matter when you choose a stack.

### Short 5: Deployment is the easy bit

**Hook on screen:** Getting an app online is no longer the interesting problem.

Most builders can deploy something now.

Then production goes wrong and they discover they don’t have the right logs, or they have thousands of metrics and no idea what any of them mean.

I’m becoming much more interested in making telemetry understandable than adding another way to deploy an app.

### Short 6: Why I’m opening DeployCrate

**Hook on screen:** I should have trusted my first instinct.

DeployCrate was originally going to be open. I let somebody outside the developer market convince me to keep it closed.

But developers want to inspect infrastructure, run it and understand it before they ask a company to pay.

So I’m opening it up and letting the product earn trust in the way developers already work.

### Short 7: I let the clanker have CRUD

**Hook on screen:** The clanker can have the CRUD screens.

Some admin table with filters, a form and a delete modal is not where I want to spend an evening.

I’ll work out the data model, permissions and backend behaviour. The agent can write the predictable interface and I’ll review it.

That split is working surprisingly well for me.

### Short 8: What MBV Labs is for

**Hook on screen:** A fuzzy software problem rarely arrives as a clean ticket.

Founders usually know something needs to be built, fixed or automated. The path is the unclear part.

I help work out what the useful version is, explain the technical tradeoffs and then build it with them.

That is what I mean when I call MBV Labs a fractional tech lead who still ships.

---

## Publishing notes

The first blog post is the main piece for the week. The thread, first LinkedIn post and several of the shorts can come from the same recording session without repeating the exact wording.

The DeployCrate posts should wait until the open source or source-available license is settled. Once it is, publish the blog and video close together so people have somewhere concrete to go.

For anything about the new admin UI, use a real screen recording. The story becomes much more credible when people can see the result and watch you inspect the generated code.

Keep the MBV Labs mention light. It belongs where the thought naturally connects to client work. Adding a call to action to every post will make the whole set feel more manufactured.
