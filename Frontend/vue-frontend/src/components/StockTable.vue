<template>
  <div class="p-6">
    <!-- inicia  header de acciones -->
    <div class="mb-4 flex flex-wrap items-center justify-between gap-4">
      <div class="flex flex-wrap items-center gap-4 flex-1 min-w-0">
        <div
          class="flex-grow min-w-[150px] sm:min-w-[200px] max-w-full sm:max-w-[400px]"
        >
          <input
            type="text"
            id="searchInput"
            v-model="searchQuery"
            placeholder="Buscar"
            class="w-full p-2 border rounded"
          />
        </div>

        <button
          @click="fetchData"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded whitespace-nowrap"
        >
          Buscar
        </button>
      </div>

      <div class="flex items-center gap-2">
        <label for="rowsPerPage">Filas:</label>
        <select
          id="rowsPerPage"
          v-model="selectedRows"
          @change="updateRowsPerPage"
          class="p-2 border rounded"
        >
          <option value="10">10</option>
          <option value="20">20</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
      </div>
    </div>
    <!-- fianaliza  header de acciones -->

    <!-- carga tabla-->
    <div v-if="stockStore.loading" class="text-center py-4">
      Cargando datos...
    </div>
    <div
      v-if="stockStore.error"
      class="bg-red-100 text-red-700 p-4 rounded mb-4"
    >
      Error: {{ stockStore.error }}
    </div>

    <div class="overflow-x-auto shadow-md rounded-lg">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th
              @click="setSort('Ticker')"
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Ticker
              <span v-if="sortColumn === 'Ticker'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              @click="setSort('Company')"
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Empresa
              <span v-if="sortColumn === 'Company'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              @click="setSort('TargetFrom')"
              class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Desde ($)
              <span v-if="sortColumn === 'TargetFrom'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              @click="setSort('TargetTo')"
              class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Hasta ($)
              <span v-if="sortColumn === 'TargetTo'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              @click="setSort('Brokerage')"
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Broker
              <span v-if="sortColumn === 'Brokerage'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              @click="setSort('Action')"
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Acción
              <span v-if="sortColumn === 'Action'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
            <th
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase"
            >
              Rating (de → a)
            </th>
            <th
              @click="setSort('Time')"
              class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase cursor-pointer"
            >
              Fecha
              <span v-if="sortColumn === 'Time'">
                {{ sortDirection === "asc" ? "▲" : "▼" }}
              </span>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr
            v-for="stock in filteredAndSortStocks"
            :key="stock.Id"
            class="hover:bg-gray-50 transition-colors"
          >
            <td class="px-4 py-3 whitespace-nowrap">{{ stock.Ticker }}</td>
            <td class="px-4 py-3">{{ stock.Company }}</td>
            <td class="px-4 py-3 text-right">
              {{ cleanPrice(stock.TargetFrom, true) }}
            </td>
            <td
              class="px-4 py-3 text-right"
              :class="getVariationClass(stock.TargetFrom, stock.TargetTo)"
            >
              {{ cleanPrice(stock.TargetTo, true) }}
              {{ getVariationIcon(stock.TargetFrom, stock.TargetTo) }}
            </td>
            <td class="px-4 py-3">{{ stock.Brokerage }}</td>
            <td class="px-4 py-3">
              <span
                :class="[
                  'px-2 py-1 rounded text-xs font-medium',
                  getActionClass(stock.Action),
                ]"
              >
                {{ translateAction(stock.Action) }}
              </span>
            </td>
            <td
              class="px-4 py-3"
              :class="getRatingClass(stock.RatingFrom, stock.RatingTo)"
            >
              {{ stock.RatingFrom }} → {{ stock.RatingTo }}
            </td>
            <td class="px-4 py-3 whitespace-nowrap">
              {{ formatDate(stock.Time) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="pagination-controls mt-4">
      <button class="p-2 border rounded" @click="prevPage">Anterior</button>
      <span class="mx-2"
        >Página {{ stockStore.currentPageKey }} de
        {{ Math.round(totalStocks / selectedRows) }}</span
      >
      <button class="p-2 border rounded" @click="nextPage">Siguiente</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch } from "vue";
import { useStockStore } from "@/stores/stockStore";

const stockStore = useStockStore();

const selectedRows = ref(10);
const searchQuery = ref("");
const searchKey = ref("");
const currentPageKey = ref(1);

const prevPage = () => {
  stockStore.currentPageKey < 2 ? 1 : stockStore.currentPageKey--;
  fetchData();
};

const nextPage = () => {
  stockStore.currentPageKey++;
  fetchData();
};

const fetchData = () => {
  searchKey.value = searchQuery.value;
  stockStore.fetchStocks(searchKey.value, selectedRows.value.toString());
};

onMounted(() => {
  fetchData();
});

const updateRowsPerPage = () => {
  stockStore.setRowsPerPage(selectedRows.value);
  fetchData();
};

const sortColumn = ref("");
const sortDirection = ref<"asc" | "desc">("asc");

const setSort = (column: string) => {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
  } else {
    sortColumn.value = column;
    sortDirection.value = "asc";
  }
};

const parseDollar = (value: string): number =>
  parseFloat(value.replace(/[^\d.-]/g, ""));

const sortedStocks = computed(() => {
  const stocks = [...stockStore.stocks];
  if (!sortColumn.value) return stocks;
  return stocks.sort((a, b) => {
    let aVal = a[sortColumn.value];
    let bVal = b[sortColumn.value];

    if (sortColumn.value === "TargetFrom" || sortColumn.value === "TargetTo") {
      aVal = parseDollar(aVal);
      bVal = parseDollar(bVal);
    }

    if (sortColumn.value === "Time") {
      aVal = new Date(aVal).getTime();
      bVal = new Date(bVal).getTime();
    }

    if (typeof aVal === "string") {
      return sortDirection.value === "asc"
        ? aVal.localeCompare(bVal)
        : bVal.localeCompare(aVal);
    } else {
      return sortDirection.value === "asc" ? aVal - bVal : bVal - aVal;
    }
  });
});

const filteredStocks = computed(() => {
  return stockStore.stocks;
});

const totalStocks = computed(() => {
  return stockStore.totalStocks;
});

const filteredAndSortStocks = computed(() => {
  const stocks = [...filteredStocks.value];
  if (!sortColumn.value) return stocks;
  return stocks.sort((a, b) => {
    let aVal = a[sortColumn.value];
    let bVal = b[sortColumn.value];

    if (sortColumn.value === "TargetFrom" || sortColumn.value === "TargetTo") {
      aVal = parseDollar(aVal);
      bVal = parseDollar(bVal);
    }

    if (sortColumn.value === "Time") {
      aVal = new Date(aVal).getTime();
      bVal = new Date(bVal).getTime();
    }

    if (typeof aVal === "string") {
      return sortDirection.value === "asc"
        ? aVal.localeCompare(bVal)
        : bVal.localeCompare(aVal);
    } else {
      return sortDirection.value === "asc" ? aVal - bVal : bVal - aVal;
    }
  });
});

const cleanPrice = (price: string, includeDollar = false) => {
  const value = parseDollar(price).toFixed(2);
  return includeDollar ? `$${value}` : value;
};

const getVariationIcon = (from: string, to: string) => {
  const numFrom = parseDollar(from);
  const numTo = parseDollar(to);
  return numTo > numFrom ? "▲" : "▼";
};

const getVariationClass = (from: string, to: string) => {
  const numFrom = parseDollar(from);
  const numTo = parseDollar(to);
  return numTo > numFrom ? "text-green-600" : "text-red-600";
};

const translateAction = (action: string) => {
  const translations: { [key: string]: string } = {
    "upgraded by": "Mejora",
    "reiterated by": "Mantiene",
    "downgraded by": "Baja",
  };
  return translations[action] || action;
};

const getActionClass = (action: string) => {
  const classes: { [key: string]: string } = {
    "upgraded by": "bg-green-100 text-green-800",
    "reiterated by": "bg-blue-100 text-blue-800",
    "downgraded by": "bg-red-100 text-red-800",
  };
  return classes[action] || "bg-gray-100 text-gray-800";
};

const getRatingClass = (from: string, to: string) => {
  if (from === to) return "text-gray-600";
  const ratingOrder = ["Sell", "Underperform", "Neutral", "Outperform", "Buy"];
  return ratingOrder.indexOf(to) > ratingOrder.indexOf(from)
    ? "text-green-700"
    : "text-red-700";
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return date.toLocaleDateString("es-ES", {
    day: "2-digit",
    month: "2-digit",
    year: "numeric",
  });
};
</script>
