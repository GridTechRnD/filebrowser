<template>
    <div class="groups-container" v-if="!layoutStore.loading">
        <div class="header">
            <h1>Groups</h1>
            <button v-if="mode == 'Listing'" class="btn btn-primary" @click="modelHandlerCreate">+ Create group</button>
            <button v-else-if="mode == 'Creating'" class="button button--flat button--red" @click="modelHandlerCreate">Return</button>
            <button v-else-if="mode == 'Editing'" class="button button--flat button--red" @click="_ => modelHandlerEdit(null)">Return</button>
        </div>
        <div v-if="mode == 'Listing'" class="table-container">
            <table class="groups-table">
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Members</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="group in groups" :key="group.groupName">
                        <td>{{ group.groupName }}</td>
                        <td>{{ group.usersIds?.length || 0 }}</td>
                        <td>
                            <button type="button" class="button button--flat" @click="modelHandlerEdit(group)">Edit</button>
                            <button type="button" class="button button--flat button--red" @click="deleteSelectedGroup(group.id)">Delete</button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <CreateGroup v-else-if="mode == 'Creating'" :modelHandlerCreate="modelHandlerCreate" />
        <EditGroup v-else :selectedGroup="selectedGroup" :modelHandlerEdit="modelHandlerEdit" />
    </div>
</template>

<script setup lang="ts">

import { inject, ref, onMounted } from "vue"
import { deleteGroup, getAllGroups } from "@/api/groups"
import { useLayoutStore } from "@/stores/layout";
import CreateGroup from "@/components/settings/CreateGroup.vue"
import EditGroup from "@/components/settings/EditGroup.vue";

const layoutStore = useLayoutStore();

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;

const mode = ref<string>("Listing")
const groups = ref<IGroup[]>([])
const selectedGroup = ref<IGroup | null>(null)

onMounted(async () => {

    try {

        layoutStore.loading = true

        groups.value = await getAllGroups()

        layoutStore.loading = false

    } catch (error) {
        layoutStore.loading = false
        $showError("Error loading groups")
    }

})

async function modelHandlerCreate() {

    if (mode.value == "Creating") {

        getAllGroups().then(response => {

            groups.value = response
            mode.value = "Listing"
        })
    }

    if (mode.value == "Listing") {
        
        mode.value = "Creating"
    }   
}

async function modelHandlerEdit( group: IGroup | null ) {

    if (mode.value == "Editing") {

        getAllGroups().then(response => {

            groups.value = response
            mode.value = "Listing"
        })
    }

    if (mode.value == "Listing") {
        
        selectedGroup.value = group
        mode.value = "Editing"
    }
}

async function deleteSelectedGroup( id: number ) {
    
    deleteGroup(id).then( _ => {
        groups.value = groups.value.filter( g => g.id !== id )
        $showSuccess("Group deleted successfully")
    }).catch( err => {
        $showError("Error deleting group", err)
    })

}
</script>

<style scoped>
.groups-container {
    padding: 1rem;
    background-color: var(--surfacePrimary);
    border-radius: 8px;
    box-shadow: 0 2px 4px var(--divider);
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}

.header h1 {
    color: var(--textSecondary);
    font-size: 1.5rem;
}

.table-container {
    overflow-x: auto;
}

.groups-table {
    width: 100%;
    border-collapse: collapse;
    background-color: var(--surfaceSecondary);
    border-radius: 8px;
    overflow: hidden;
}

.groups-table th, .groups-table td {
    padding: 0.75rem;
    text-align: left;
    color: var(--textPrimary);
}

.groups-table th {
    background-color: var(--surfacePrimary);
    font-weight: bold;
}

.groups-table tr:nth-child(even) {
    background-color: var(--background);
}

.groups-table tr:hover {
    background-color: var(--hover);
}

.btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.btn-primary {
    background-color: var(--blue);
    color: white;
}

.btn-primary:hover {
    background-color: var(--dark-blue);
}

.btn-secondary {
    background-color: var(--red);
    color: white;
}

.btn-secondary:hover {
    background-color: var(--dark-red);
}

.btn-edit {
    background-color: var(--icon-blue);
    color: white;
    margin-right: 0.5rem;
}

.btn-edit:hover {
    background-color: var(--dark-blue);
}

.btn-delete {
    background-color: var(--icon-red);
    color: white;
}

.btn-delete:hover {
    background-color: var(--dark-red);
}
</style>