# Creating or updating entity files
When a new chapter is written and approved, if existing entities have developed, or new entities have been created, create or update entity files.

### What Needs to Happen Post-Chapter
```
Task                                        Complexity      Needs Context
1. Extract what entities appeared           Low             Chapter text only
2. Detect what's new vs existing            Low (code)      Entity index
3. Diff what changed for existing entities  Medium          Chapter + existing entity file
4. Generate updates for existing entities   Low per entity  Existing file + diff5. Generate new entity filesMedium per entityChapter + templates
```

### Recommended Approach: Extraction Agent + Parallel Haiku Calls

```
New Chapter
    │
    ▼
┌─────────────────────────────────────────────┐
│  Agent 6: Entity Extractor (Haiku)          │
│  - One call                                 │
│  - Lists all entities mentioned             │
│  - Notes what happened to each              │
│  - Flags new vs existing                    │
└─────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────┐
│  Code: Diff & Dispatch                      │
│  - Match extracted entities to existing     │
│  - Determine: UPDATE, CREATE, or SKIP       │
│  - Queue generation tasks                   │
└─────────────────────────────────────────────┘
    │
    ├──▶ Haiku: Update char_001 (parallel)
    ├──▶ Haiku: Update loc_001 (parallel)
    ├──▶ Haiku: Create char_005 (parallel)
    └──▶ Haiku: Create obj_004 (parallel)
```

### Why Haiku for everything?
These are structured extraction/generation tasks, not creative
Entity updates are formulaic: "add this event, update this field"

### Pipeline Flow

```
Chapter 16 Written & Approved
           │
           ▼
    ┌──────────────────────┐
    │ Agent 6: Extractor   │  1 Haiku call
    │ (Haiku)              │  ~$0.002
    └──────────────────────┘
           │
           ▼ extraction.json
           │
    ┌──────────────────────┐
    │ Code: Router         │  No API call
    │ - Match to index     │
    │ - Filter SKIP        │
    │ - Queue tasks        │
    └──────────────────────┘
           │
           ├─────────────────────────────────────────┐
           │                                         │
           ▼                                         ▼
    ┌─────────────────┐                    ┌─────────────────┐
    │ UPDATE tasks    │                    │ CREATE tasks    │
    │ (parallel)      │                    │ (parallel)      │
    └─────────────────┘                    └─────────────────┘
           │                                         │
    ┌──────┴──────┐                          ┌──────┴──────┐
    ▼             ▼                          ▼             ▼
 Agent 7      Agent 7                     Agent 8      Agent 8
 Update       Update                      Create       Create
 char_001     char_002                    obj_004      (if any)
 ~$0.002      ~$0.002                     ~$0.003
```

### Cost Per Chapter
```
Scenario                        Haiku Calls     Est. Cost
Typical (2 updates, 0 new)      3               ~$0.006
Active (3 updates, 1 new)       5               ~$0.011
Big chapter (4 updates, 2 new)  7               ~$0.016
```

### Entity file updating
Entity Updater sees extraction:
`Chapter 16: Trust destroyed, voice changed, new physical tell`

Entity Updater does:
1. Overwrites personality.core_trait.current → "hardened"
2. Appends to personality.core_trait.evolution → new entry
3. Overwrites voice.current → new speech pattern
4. Appends to voice.evolution → new chapter range entry
5. Adds betrayal to physical_tells.current
6. Appends to physical_tells.evolution → notes when/why added
7. Appends to key_events → chapter 16 event
8. Updates last_appeared_chapter → 16

### Go code for implementation:
``` golang
// After chapter is written and approved
func (p *Pipeline) UpdateEntities(ctx context.Context, chapter Chapter) error {
    // Step 1: Extract (single Haiku call)
    extraction, err := p.entityExtractor.Run(ctx, chapter)
    if err != nil {
        return err
    }
    
    // Step 2: Route
    var updates []EntityTask
    var creates []EntityTask
    
    for _, e := range extraction.Entities {
        if e.Priority == "skip" {
            continue
        }
        if e.Status == "existing" {
            updates = append(updates, EntityTask{
                Type:       "update",
                EntityID:   e.EntityID,
                Extraction: e,
            })
        } else {
            creates = append(creates, EntityTask{
                Type:       "create",
                EntityType: e.Type,
                Extraction: e,
            })
        }
    }
    
    // Step 3: Execute in parallel
    var wg sync.WaitGroup
    errCh := make(chan error, len(updates)+len(creates))
    
    for _, task := range updates {
        wg.Add(1)
        go func(t EntityTask) {
            defer wg.Done()
            if err := p.entityUpdater.Run(ctx, t); err != nil {
                errCh <- err
            }
        }(task)
    }
    
    for _, task := range creates {
        wg.Add(1)
        go func(t EntityTask) {
            defer wg.Done()
            if err := p.entityCreator.Run(ctx, t); err != nil {
                errCh <- err
            }
        }(task)
    }
    
    wg.Wait()
    close(errCh)
    
    // Collect errors
    var errs []error
    for err := range errCh {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("entity updates failed: %v", errs)
    }
    
    return nil
}
```