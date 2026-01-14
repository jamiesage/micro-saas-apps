# World Tracking System

## Overview

The world tracking system has two layers:

| File | Purpose | Updates |
|------|---------|---------|
| `world_map.json` | **Topology** — locations, connections, terrain, distances | When locations discovered or connections change |
| `world_state.json` | **Current state** — where everyone/everything IS right now | After every chapter |

## How They Work Together

```
world_map.json (static geography)
├── Locations: Where places ARE (coordinates, terrain, connections)
├── Connections: How places link (routes, travel times, hazards)
└── Regions: Groupings for narrative/map purposes

world_state.json (dynamic positions)
├── character_positions: Where each character is NOW
├── creature_positions: Where creatures are NOW  
├── object_positions: Where objects are NOW
├── recent_movements: What moved last few chapters
├── location_occupancy: Quick "who's here?" lookup
└── travel_in_progress: Active journeys
```

## Pipeline Integration
```
Chapter Written
     │
     ▼
Entity Extractor ─────────────────┐
     │                            │
     │ entities_extracted         │ movements_detected
     │                            │
     ▼                            ▼
Entity Updater (parallel)    Position Updater
     │                            │
     ▼                            ▼
Entity files updated         world_state.json updated
(with location.history)      (current positions)
```

## For Claude Agents

### Story Planner reads:
```json
// "Where can they go from here?"
world_map.connections.filter(c => c.from === current_location)

// "How long would that take?"
world_map.distance_matrix[from][to]

// "Who else is nearby?"
world_state.location_occupancy[current_location]

// "Are any characters about to meet?"
world_state.proximity_alerts
```

### Story Writer reads:
```json
// "How to describe this movement?"
world_map.cardinal_directions["loc_001_to_loc_003"]
// → "east through the forest depths"

// "What terrain are they in?"
world_map.locations.find(l => l.id === current).terrain
// → "dense_forest"

// "Who's following them?"
world_state.creature_positions.filter(c => c.following === "char_001")
// → Thornback
```

### Entity Extractor flags:
```json
// "Did anyone move this chapter?"
{
  "movement_detected": [
    {"entity": "char_001", "from": "standing stones", "to": "Kael's camp"}
  ]
}
```

### Position Updater (new small agent) updates:
```json
// world_state.json after chapter 17
{
  "character_positions": [
    {
      "entity_id": "char_001",
      "current_location": {
        "sub_location": "Confrontation site, eastern Thornwood"
        // ... updated fields
      }
    }
  ]
}
```

## Coordinate System

```
                    NORTH (y+)
                       │
                       │
    WEST (x-)  ────────┼──────── EAST (x+)
                       │
                       │
                    SOUTH (y-)

Grid unit = ~1 day travel on clear terrain
Origin (0,0) = Gallows Crossing

Example positions:
  Gallows Crossing:  (0, 0)
  Thornwood center:  (2, 0)  
  Eastern edge:      (4, 0)
  Mountain Pass:     (5, 1)
  Hidden Vale:       (7, 0)
```

## Travel Time Calculation

```
Base time = distance_matrix[from][to]
Actual time = base_time / terrain.travel_modifier

Example:
  Thornwood crossing = 14 days base
  Dense forest modifier = 0.3
  But that's already factored in the matrix
  
  If route blocked or character injured:
  Actual time = base_time * complication_modifier
```

## Position Types

### Characters
```json
{
  "movement_status": "traveling | stationary | unknown | fleeing",
  "heading": "north | east | south | west | none | unknown",
  "destination": "loc_id or null",
  "coordinates_approx": {"x": 3.5, "y": 0.2}
}
```

### Creatures
```json
{
  "movement_status": "territorial | following | migrating | stationary | unknown",
  "following": "char_id or null",
  "territory": ["loc_ids..."],
  "distribution": "scattered | concentrated | settlement_and_patrols"
}
```

### Objects
```json
{
  "position_type": "carried | placed | hidden | lost | given_away | destroyed",
  "carrier_id": "entity carrying it, or null",
  "location_id": "if placed/hidden, where"
}
```

## Update Flow

```
Chapter Written
     │
     ▼
Entity Extractor (existing)
     │ Detects: "Mira moved to the river"
     │ Detects: "She dropped the compass"
     ▼
Position Updater (small Haiku agent)
     │ Input: extraction + current world_state
     │ Output: updated world_state.json
     ▼
world_state.json updated
     │
     ▼
Entity files updated (location.history appended)
```

## Position Updater Agent

A small dedicated agent for world state updates:

**Input:**
- Movement extractions from Entity Extractor
- Current world_state.json
- Current world_map.json (for validation)

**Tasks:**
1. Update character_positions
2. Update creature_positions (if any moved)
3. Update object_positions (if any changed hands)
4. Append to recent_movements
5. Update location_occupancy
6. Update travel_in_progress
7. Check for new proximity_alerts

**Output:** Updated world_state.json

## Map Generation (Future)

The coordinate system enables automated map generation:

```javascript
// Pseudocode for map visualization
locations.forEach(loc => {
  plot(loc.position.x, loc.position.y, {
    icon: getIcon(loc.type),
    label: loc.name,
    discovered: loc.discovered
  });
});

connections.forEach(conn => {
  drawLine(
    locations[conn.from].position,
    locations[conn.to].position,
    {
      style: getStyle(conn.type),
      dashed: conn.status === 'unknown'
    }
  );
});

// Character positions overlay
characterPositions.forEach(char => {
  plotCharacter(char.coordinates_approx, {
    icon: getCharacterIcon(char.entity_id),
    label: char.name
  });
});
```

## Spatial Queries

Common queries the system supports:

| Query | How to Answer |
|-------|---------------|
| "Where is Mira?" | `world_state.character_positions.find(c => c.entity_id === "char_001")` |
| "Who's in the Thornwood?" | `world_state.location_occupancy["loc_001"]` |
| "How far to the mountains?" | `world_map.distance_matrix["loc_001"]["loc_004"]` |
| "What's between here and there?" | `world_map.connections.find(c => c.from === here && c.to === there)` |
| "Are any characters near each other?" | `world_state.proximity_alerts` |
| "What route is she taking?" | `world_state.travel_in_progress.find(t => t.entity_id === "char_001")` |
| "Has she been here before?" | Entity file: `location.history` |

## Consistency Rules

1. **Entity files are source of truth for history** — location.history in each entity tracks where they've been
2. **world_state is source of truth for NOW** — current positions live here, not duplicated
3. **world_map is source of truth for geography** — distances, connections, terrain don't change (unless story changes them)
4. **Coordinates are approximate** — sub_location text is more important for prose
5. **Unknown is valid** — characters can have unknown positions (Kael, Brennan)