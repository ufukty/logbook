import { Task } from "./Task.js";

export function ChronologicalDocumentOverview(props) {
    // Properties that are fetched from server
    this.status = props.status; // example: "8557d156-3d00-4836-8323-a9bdd586719a"
    this.incident_id = props.incident_id; // example: "61bbc44a-c61c-4d49-8804-486181081fa7"
    this.error_hint = props.error_hint; // example: "999c060e-d853-4271-b529-42c2655a4aae"
    this.resource = props.resource.map((resourceItem) => {
        return new Task(resourceItem);
    }); // example: "Update redis/tf file according to prod.tfvars file"
}
