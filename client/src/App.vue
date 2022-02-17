<template>
  <div class="main">
    <ag-grid-vue
      class="ag-theme-alpine grid-container"
      :columnDefs="columnDefs"
      @grid-ready="onGridReady"
      :rowSelection="rowSelection"
      :rowData="rowData"
      @selection-changed="onSelectionChanged"
      @cell-value-changed="onCellValueChanged"
    ></ag-grid-vue>
    <div class="container">
      <span v-if="selectedTodo">
        <span>delete selected todo </span>
        <button @click="removeTodo">Delete</button>
      </span>
      <add-todo @create-todo="createTodo" />
    </div>
  </div>
</template>
<script>
import "ag-grid-community/dist/styles/ag-grid.css";
import "ag-grid-community/dist/styles/ag-theme-alpine.css";
import { AgGridVue } from "ag-grid-vue3";
import axios from "axios";

import AddTodo from "../components/AddTodo.vue";

export default {
  components: {
    "ag-grid-vue": AgGridVue,
    AddTodo,
  },
  data: function () {
    return {
      endpoint: "http://localhost:5000",
      columnDefs: [
        {
          field: "title",
          editable: true,
        },
        {
          field: "body",
          editable: true,
        },
        {
          field: "completed",
          editable: true,
          flex: 1,
        },
      ],

      gridApi: null,
      columnApi: null,
      rowSelection: null,
      rowData: null,
      selectedTodo: null,
      todoForm: {
        title: "",
        body: "",
        edit: false,
        completed: false,
      },
    };
  },
  created() {
    this.rowSelection = "single";
  },
  methods: {
    onSelectionChanged() {
      const selectedTodo = this.gridApi.getSelectedRows();
      this.selectedTodo = selectedTodo[0];
    },
    async onGridReady(params) {
      this.gridApi = params.api;
      this.gridColumnApi = params.columnApi;

      await axios
        .get(this.endpoint + "/todos")
        .then((res) => (this.rowData = res.data));
    },
    onCellValueChanged(params) {
      const { data, newValue, colDef, node, oldValue } = params;

      const updatedTodo = JSON.stringify({
        ...this.rowData[node.rowIndex],
        [colDef.field]: newValue,
      });

      axios
        .patch(this.endpoint + "/editTodo/" + data.id, updatedTodo)
        .then((res) => {
          res.status === 200
            ? this.filterTodo(colDef.field, newValue, data.id)
            : this.filterTodo(colDef.field, oldValue, data.id);
        })
        .catch((error) => {
          filterTodo(colDef.field, oldValue, data.id);
          alert(error);
        });
    },
    filterTodo(key, value, id) {
      this.rowData = this.rowData.map((todo) => {
        if (todo.id == id) {
          return { ...todo, [key]: value };
        }
        return todo;
      });
    },
    removeTodo() {
      axios
        .delete(this.endpoint + "/todos/" + this.selectedTodo.id)
        .then((res) => {
          if (res.status === 200) {
            this.rowData = this.rowData.filter(
              (todo) => todo.id !== res.data.id
            );
            this.selectedTodo = null;
          }
        })
        .catch((error) => {
          alert(error);
        });
    },
    createTodo(e) {
      axios
        .post(this.endpoint + "/todo", JSON.stringify(e))
        .then((res) => {
          this.rowData = [...this.rowData, res.data];
        })
        .catch((error) => {
          alert(error);
        });
    },
  },
};
</script>

<style lang="scss">
@import "~ag-grid-community/dist/styles/ag-grid.css";
@import "~ag-grid-community/dist/styles/ag-theme-alpine.css";

.main {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

.grid-container {
  width: 510px;
  overflow-x: hidden !important;
  height: 500px;
}

.container {
  margin-top: 15px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}
</style>
