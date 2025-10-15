<template>
  <div class="settings-modal">
    <!-- Nome do grupo -->
    <div class="form-group">
      <label for="group-name" class="form-label">Group name:</label>
      <input
        id="group-name"
        v-model="groupName"
        class="input"
        placeholder="Ex: Back Office"
      />
    </div>

    <!-- Lista de usuários selecionados -->
    <div class="form-section">
      <h3 class="section-title">Users</h3>
      <table class="listing-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in selectedUsers" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td><button type="button" style="background-color: red; color: white; border: 0px; width: 20px;" @click="handleRemoveUser(user.id)">-</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Adicionar usuário -->
    <div class="form-group user-select" style="margin: 20px;">
      <select v-model="selectedUserId" class="input">
        <option disabled value="">-- Select a user --</option>
        <option
          v-for="user in availableUsers"
          :key="user.id"
          :value="user.id"
        >
          {{ user.username }}
        </option>
      </select>
      <button class="button button--flat" @click="handleAddUser" :disabled="!selectedUserId">
        Add
      </button>
    </div>

    <!-- Regras -->
    <div class="form-section">
      <h3 class="section-title">Rules</h3>
      <Rules v-model:rules="rules" />
    </div>

    <!-- Ações -->
    <div class="form-actions">
      <button class="button success" @click="handleEditGroup">
        Save
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">

import { inject, ref, onMounted } from "vue";
import { groups as gApi } from "@/api";
import { users as uApi } from "@/api";
import { useLayoutStore } from "@/stores/layout";
import Rules from "./Rules.vue";

interface SelectionUser {
    id: number;
    username: string;
}

const props = defineProps([
    "selectedGroup",
    "modelHandlerEdit"
])

const layoutStore = useLayoutStore();

const $showError = inject<IToastError>("$showError")!;
const $showSuccess = inject<IToastSuccess>("$showSuccess")!;

const groupName = ref<string>(props.selectedGroup.groupName);
const availableUsers = ref<SelectionUser[]>([]);
const selectedUsers = ref<SelectionUser[]>([]);
const selectedUserId = ref<number | null>(null);
const rules = ref<IRule[]>([...props.selectedGroup.groupRules])

// Buscar usuários ao montar
onMounted(async () => {
  try {

    const response = await uApi.getAll();

    availableUsers.value = response
      .filter((user: any) => !user.perm.admin && !props.selectedGroup.usersIds?.includes(user.id))
      .map((user: any) => ({ id: user.id, username: user.username }));

    selectedUsers.value = response
      .filter((user: any) => props.selectedGroup.usersIds?.includes(user.id))
      .map((user: any) => ({ id: user.id, username: user.username }));  
    
  } catch (err) {
    alert(err);
  }
});

function handleAddUser() {
  const user = availableUsers.value.find((u) => u.id === selectedUserId.value);
  if (user) {
    selectedUsers.value.push(user);
    availableUsers.value = availableUsers.value.filter((u) => u.id !== user.id);
    selectedUserId.value = null;
  }
}

function handleRemoveUser(userId: number) {
  const user = selectedUsers.value.find((u) => u.id === userId);
  if (user) {
    availableUsers.value.push(user);
    selectedUsers.value = selectedUsers.value.filter((u) => u.id !== userId);
  }
}

function handleEditGroup() {
  gApi.updateGroup(
    props.selectedGroup.id,
    groupName.value.trim(),
    rules.value,
    selectedUsers.value.map((u) => u.id)
  ).then( _ => {
    props.modelHandlerEdit().then( () => {

      $showSuccess("Group edited successfully!");

      groupName.value = "";
      selectedUsers.value = [];
      rules.value = [
          {
              allow: true,
              regex: false,
              path: "",
              regexp: {
              raw: "",
              },
          },
      ];
    });
        
        
    }).catch( err => {
        $showError("Failed to create group.", err);
    });
}
</script>

<style scoped>
.user-select {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}
</style>
