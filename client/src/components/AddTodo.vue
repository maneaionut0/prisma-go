<template>
  <input type="text" v-model="todoForm.title" placeholder="title " />
  <input type="text" v-model="todoForm.body" placeholder="body " />
  <button @click="createTodo">Add Todo</button>
</template>
<script>
export default {
  name: "add-todo",
  emits: ["createTodo"],
  data() {
    return {
      todoForm: {
        title: "",
        body: "",
      },
    };
  },
  methods: {
    createTodo() {
      this.$emit("createTodo", this.todoForm);
      this.todoForm = {
        title: "",
        body: "",
      };
    },
  },
};
</script>
<style></style>
