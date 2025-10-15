
import { fetchJSON, fetchURL } from "./utils";

export async function createGroup(name: string, rules: IRule[], users: number[]) {
    return fetchURL(`/api/groups`, {
        method: "POST",
        body: JSON.stringify({
            what: "groups",
            witch: [],
            data: {
                groupName: name,
                usersIds: users,
                groupRules: rules
            }
        })
    });
}

export async function getAllGroups() {
    return fetchJSON<IGroup[]>(`/api/groups`, {});
}

export async function updateGroup(id: number, name: string, rules: IRule[], users: number[]) {
    return fetchURL(`/api/groups/${id}`, {
        method: "PUT",
        body: JSON.stringify({
            what: "groups",
            which: ["all"],
            data: {
                ID: id,
                groupName: name,
                usersIds: users,
                groupRules: rules
            }
        })
    });
}

export async function deleteGroup(id: number) {
    return fetchURL(`/api/groups/${id}`, {
        method: "DELETE",
    });
}

