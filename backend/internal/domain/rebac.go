package domain

type Entity struct {
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	EntityTable string `json:"entity_table"`
}

// Relationship model defines the relationship between entities
//
// There some relation predefined in the system:
//   - owner (the subject is the owner of the entity;
//     an owner has full access to the entity)
//   - member (the subject is a member of the entity;
//     a member has access to the entity based on the entity's
//     permissions)
//   - parent (the subject is the parent of the entity, the entity inherits relations and permissions from the parent)
//   - viewer (the subject can view the entity)
//   - editor (the subject can edit the entity)
//
// Example:
// - document:1#owner@user:1 (user with ID 1 is the owner of the document with ID 1)
// - document:1#viewer@group:1#member (members of group with ID 1 is a viewer of the document with ID 1)
type RelationTuple struct {
	EntityType      string `json:"entity_type"`
	EntityID        string `json:"entity_id"`
	Relation        string `json:"relation"`
	SubjectType     string `json:"subject_type"`
	SubjectID       string `json:"subject_id"`
	SubjectRelation string `json:"subject_relation"`
}

// Example:
// - document#owner@user
// - document#viewer@group#member
// - document#view@owner
// - document#edit@owner
// - document#delete@owner
// - document#view@project#member

type RelationDefinition struct {
	EntityType      string `json:"entity_type"`
	RelationType    string `json:"relation_type"`
	SubjectType     string `json:"subject_type"`
	SubjectRelation string `json:"subject_relation"`
}

type Permission struct {
	Relation string `json:"relation"`
	Action   string `json:"action"`
}
