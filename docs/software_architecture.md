**About arc42**

arc42, the Template for documentation of software and system
architecture.

By Dr. Gernot Starke, Dr. Peter Hruschka and contributors.

Template Revision: 7.0 EN (based on asciidoc), January 2017

© We acknowledge that this document uses material from the arc 42
architecture template, <http://www.arc42.de>. Created by Dr. Peter
Hruschka & Dr. Gernot Starke.

Introduction and Goals
======================
tournaments shall organize tcg-tournaments with an HTTP-API using Hyper-Items.

Requirements Overview
---------------------
The overall goal of tournaments is to provide an API for managing:
* Decks
* Players
* Seasons

### Decks
A Decks here is either a combination of a name and a link to a Deck-Builder
or github.com/cognicraft/mtg dependency for parsing decks
* Available Decks
* Adding Decks
* Deleting Decks

### Players
A Player is a profile which contains various information like:
* Profile picture
* Name
* Winrate
* Favourite Decks
* Seasons participation

### Seasons
A Season is a time restricted event in which players may participate and play matches according
to the Season's format. A Season will require:
* Format
* Players and their chosen Decks
And will show various stats like:
* Standings
  * If the Season is still ongoing: shows current standings
  * If the Season is finished: shows summary of Season and final standings
* Brackets, if the current formats allows it 

Quality Goals {#_quality_goals}
-------------

+-------------+---------------------------+---------------------------+
| Priority    | Quality Goal              | Scenario                  |
+=============+===========================+===========================+
| *\<Role-1\> | *\<Contact-1\>*           | *\<Expectation-1\>*       |
| *           |                           |                           |
+-------------+---------------------------+---------------------------+
| *\<Role-2\> | *\<Contact-2\>*           | *\<Expectation-2\>*       |
| *           |                           |                           |
+-------------+---------------------------+---------------------------+

Stakeholders {#_stakeholders}
------------

+-------------+---------------------------+---------------------------+
| Role/Name   | Contact                   | Expectations              |
+=============+===========================+===========================+
| *\<Role-1\> | *\<Contact-1\>*           | *\<Expectation-1\>*       |
| *           |                           |                           |
+-------------+---------------------------+---------------------------+
| *\<Role-2\> | *\<Contact-2\>*           | *\<Expectation-2\>*       |
| *           |                           |                           |
+-------------+---------------------------+---------------------------+

Architecture Constraints
========================

* runnable from the command line
* platform-independent and should run on the major operating systems

System Scope and Context {#section-system-scope-and-context}
========================

Business Context {#_business_context}
----------------

**\<Diagram or Table\>**

**\<optionally: Explanation of external domain interfaces\>**

Technical Context {#_technical_context}
-----------------

**\<Diagram or Table\>**

**\<optionally: Explanation of technical interfaces\>**

**\<Mapping Input/Output to Channels\>**

Solution Strategy {#section-solution-strategy}
=================

Building Block View {#section-building-block-view}
===================

Whitebox Overall System {#_whitebox_overall_system}
-----------------------

***\<Overview Diagram\>***

Motivation

:   *\<text explanation\>*

Contained Building Blocks

:   *\<Description of contained building block (black boxes)\>*

Important Interfaces

:   *\<Description of important interfaces\>*

### \<Name black box 1\> {#__name_black_box_1}

*\<Purpose/Responsibility\>*

*\<Interface(s)\>*

*\<(Optional) Quality/Performance Characteristics\>*

*\<(Optional) Directory/File Location\>*

*\<(Optional) Fulfilled Requirements\>*

*\<(optional) Open Issues/Problems/Risks\>*

### \<Name black box 2\> {#__name_black_box_2}

*\<black box template\>*

### \<Name black box n\> {#__name_black_box_n}

*\<black box template\>*

### \<Name interface 1\> {#__name_interface_1}

...

### \<Name interface m\> {#__name_interface_m}

Level 2 {#_level_2}
-------

### White Box *\<building block 1\>* {#_white_box_emphasis_building_block_1_emphasis}

*\<white box template\>*

### White Box *\<building block 2\>* {#_white_box_emphasis_building_block_2_emphasis}

*\<white box template\>*

...

### White Box *\<building block m\>* {#_white_box_emphasis_building_block_m_emphasis}

*\<white box template\>*

Level 3 {#_level_3}
-------

### White Box \<\_building block x.1\_\> {#_white_box_building_block_x_1}

*\<white box template\>*

### White Box \<\_building block x.2\_\> {#_white_box_building_block_x_2}

*\<white box template\>*

### White Box \<\_building block y.1\_\> {#_white_box_building_block_y_1}

*\<white box template\>*

Runtime View {#section-runtime-view}
============

\<Runtime Scenario 1\> {#__runtime_scenario_1}
----------------------

-   *\<insert runtime diagram or textual description of the scenario\>*

-   *\<insert description of the notable aspects of the interactions
    between the building block instances depicted in this diagram.\>*

\<Runtime Scenario 2\> {#__runtime_scenario_2}
----------------------

... {#_}
---

\<Runtime Scenario n\> {#__runtime_scenario_n}
----------------------

Deployment View {#section-deployment-view}
===============

Infrastructure Level 1 {#_infrastructure_level_1}
----------------------

***\<Overview Diagram\>***

Motivation

:   *\<explanation in text form\>*

Quality and/or Performance Features

:   *\<explanation in text form\>*

Mapping of Building Blocks to Infrastructure

:   *\<description of the mapping\>*

Infrastructure Level 2 {#_infrastructure_level_2}
----------------------

### *\<Infrastructure Element 1\>* {#__emphasis_infrastructure_element_1_emphasis}

*\<diagram + explanation\>*

### *\<Infrastructure Element 2\>* {#__emphasis_infrastructure_element_2_emphasis}

*\<diagram + explanation\>*

...

### *\<Infrastructure Element n\>* {#__emphasis_infrastructure_element_n_emphasis}

*\<diagram + explanation\>*

Cross-cutting Concepts {#section-concepts}
======================

*\<Concept 1\>* {#__emphasis_concept_1_emphasis}
---------------

*\<explanation\>*

*\<Concept 2\>* {#__emphasis_concept_2_emphasis}
---------------

*\<explanation\>*

...

*\<Concept n\>* {#__emphasis_concept_n_emphasis}
---------------

*\<explanation\>*

Design Decisions {#section-design-decisions}
================

Quality Requirements {#section-quality-scenarios}
====================

Quality Tree {#_quality_tree}
------------

Quality Scenarios {#_quality_scenarios}
-----------------

Risks and Technical Debts {#section-technical-risks}
=========================

Glossary {#section-glossary}
========

+-----------------------------------+-----------------------------------+
| Term                              | Definition                        |
+===================================+===================================+
| \<Term-1\>                        | \<definition-1\>                  |
+-----------------------------------+-----------------------------------+
| \<Term-2\>                        | \<definition-2\>                  |
+-----------------------------------+-----------------------------------+
