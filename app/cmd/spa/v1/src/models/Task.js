export function TaskModel(props) {
    // Properties that are fetched from server
    this.taskId = props.task_id; // example: "8557d156-3d00-4836-8323-a9bdd586719a"
    this.documentId = props.document_id; // example: "61bbc44a-c61c-4d49-8804-486181081fa7"
    this.parentId = props.parent_id; // example: "999c060e-d853-4271-b529-42c2655a4aae"
    this.content = props.content; // example: "Update redis/tf file according to prod.tfvars file"
    this.degree = props.degree; // example: 1
    this.depth = props.depth; // example: 4
    this.createdAt = props.created_at; // example: "2022-01-27T01:39:51.320386Z"
    this.completedAt = props.completed_at; // example: "2022-02-17"
    this.readyToPickUp = props.ready_to_pick_up; // example: true

    // Properties that are created and used by frontend
    this.effectiveDepth = 0;
}