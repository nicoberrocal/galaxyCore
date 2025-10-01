# Galaxy Core: Feature Overview

This document provides a brief explanation of every major feature in Galaxy Core.

## 1. Core Gameplay

Galaxy Core is a slow-paced, persistent browser-based MMORTS. The core loop involves exploring the galaxy, expanding your empire by colonizing star systems, exploiting resources, and engaging in strategic combat with other players.

## 2. Territory & Colonization

- **Systems**: Star systems are the primary military and strategic points of interest. They can be colonized and defended with fleets.
- **Planets**: Located within systems, planets are the economic backbone of an empire. They house buildings for resource and energy production.

## 3. Economy & Buildings

- **Resources**: The game features five core resources: Metal, Crystal, Hydrogen, Plasma, and Energy.
- **Buildings**: Players can construct various buildings on planets:
  - **Resource Production**: Metal Mines, Crystal Mines.
  - **Energy Production**: Solar Farms, Wind Farms, Hydro-Electric Dams, etc., with efficiency depending on the planet type.
  - **Special Facilities**: Shipyards for producing ships, and advanced structures like Particle Accelerators and Fusion Reactors.

## 4. Fleet & Ship Customization

- **Ship Blueprints**: There are 6 core ship types, each with a specific role: Drone, Scout, Fighter, Bomber, Carrier, and Destroyer.
- **Role Modes**: Every ship can switch between four modes (Tactical, Economic, Recon, Scientific) to gain situational bonuses and abilities.
- **Abilities**: Ships come with a unique set of passive, active, and conditional abilities. Key abilities include:
  - **Recon & Intel**: `LongRangeSensors` (passive vision boost), `Ping` (active target marking), `DecoyBeacon` (creates phantom contacts).
  - **Mobility**: `LightSpeed` (enables FTL warp), `WarpStabilizer` (improves warp safety), `RapidRedeploy` (reduces warp charge time).
  - **Offense**: `AlphaStrike` (first-volley damage boost), `Overload` (risky damage spike), `FocusFire` (concentrates damage), `PhaseLance` (partially ignores shields).
  - **Defense**: `PointDefenseScreen` (mitigates laser fire), `EvasiveManeuvers` (temporary evasion boost), `InterdictorPulse` (blocks enemy warp).
  - **Siege & Support**: `SiegePayload` (bonus vs. structures), `TargetingUplink` (accuracy/crit boost), `BunkerBuster` (heavy damage to fortifications).
  - **Economy**: `ResourceHarvester` (enables mining), `SelfRepair` (out-of-combat regen), `CloakWhileAnchored` (stealth while mining).

## 5. Cosmic Gem System

This system allows for deep customization of ships.

- **Gem Sockets**: Ships have sockets to equip powerful gems.
- **Gem Families & Tiers**: Gems belong to one of 8 families (e.g., Laser, Nuclear, Warp) and range from Tier 1 to 5.
- **Synthesis**: Players can upgrade gems by combining three of the same type or experiment with high-risk, high-reward synthesis by mixing different families to create unique **Hybrid** and **Relic** gems.
- **GemWords**: Socketing specific sequences of gems (e.g., `[Laser, Laser, Nuclear]`) unlocks powerful `GemWord` bonuses and grants new abilities.
- **Ultimate Synthesis**: The most powerful item is the **Singularity Core**, a legendary relic created by synthesizing **three distinct Tier 5 Hybrid Gems**. It grants massive bonuses to all damage and reduces ability cooldowns.

## 6. Combat Mechanics

- **Ship Stacks & HP Buckets**: Ships are grouped in stacks. Damage is applied to HP "buckets" (a group of ships of the same type), not individual units, making large-scale combat efficient.
- **Formations & Counters**: Players can arrange ships into tactical formations. Each has strengths and weaknesses in a rock-paper-scissors dynamic:
  - **Line**: Balanced offense/defense. Strong against `Vanguard`.
  - **Box**: Heavy all-around defense, but slow. Strong against `Line`.
  - **Vanguard**: Aggressive, forward-focused. Strong against `Box` but weak to flanking.
  - **Skirmish**: Mobile and evasive. Strong against `Vanguard`.
  - **Echelon**: Staggered lines. Strong against `Skirmish`.
- **Combat Resolution**: Combat is turn-based. The attacker fires first, followed by the defender's counter-attack. The attacking stack then enters a one-hour cooldown.
- **Damage & Shields**: There are three attack types: **Laser**, **Nuclear**, and **Antimatter**. Ships have corresponding shield values, creating a tactical triangle where choosing the right attack type is critical.

## 7. Movement & Actions

- **Action Queues**: All significant player actions, such as construction, research, and fleet movement, are handled through a time-based queue system, reinforcing the game's slow-paced nature.
- **Warp Travel**: Fleets travel between star systems using FTL/warp, which involves charge times and can be disrupted by enemy interdiction abilities.
