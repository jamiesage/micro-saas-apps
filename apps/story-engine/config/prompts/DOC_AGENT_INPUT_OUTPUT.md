# The I/O for each agent

## Code for context builder files:
Some models need more context than others. Here is a snippet on how data would be taken from entity files to form suitable contexts for different models.
```go
// For writing agents, extract only current state
func (e *Entity) CurrentState() EntityCurrent {
    return EntityCurrent{
        Name:        e.Name,
        Status:      e.Status.Current,
        Location:    e.Location.Current,
        CoreTrait:   e.Personality.CoreTrait.Current,
        Wants:       e.Personality.Wants.Current,
        Voice:       e.Voice.Current,
        Relationships: extractCurrentRelationships(e.Relationships),
        // No evolution data—just what's needed now
    }
}

// For planning/analysis, full entity available
func (e *Entity) FullHistory() Entity {
    return *e // Everything, including evolution
}
```

## Story Planner

### Input / Context
Full context files for any relevant characters. Give access to read-only any of a selection of files.

The full file with evolution tracking lets the planner:
- See how a character changed: "She went from determined → guarded → hardened"
- Track relationship arcs: "Trust built over chapters 5-9, destroyed in 16"
- Verify consistency: "Her voice change matches her arc beats"
- Plan future beats: "She's at her lowest—time for turning point"

## Story Writer

### Input / Context
Only give the story writer current data, doesn't need to know about past / future, handled by story planner.

Use a context builder to extract .current values from fields in entity files.

```json
{
  "name": "Mira Thorne",
  "personality": {
    "core_trait": "hardened",           // Just current
    "wants": {
      "external": "Find her brother",    // Just current
      "internal": "Make Kael answer"
    }
  },
  "voice": {
    "speech_pattern": "Minimal. Lets silence do the work."  // Just current
  },
  "relationships": [{
    "entity_name": "Kael Vorn",
    "type": "adversary",                 // Just current
    "status": "trust_destroyed"
  }]
}
```
