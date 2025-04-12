<template>
  <div class="min-h-screen bg-gray-100 text-gray-900 text-center">
    <h3 class="text-2xl font-bold p-6">Acciones recomendadas</h3>
    <Carousel />
    <StockTable />
  </div>
</template>

<script setup lang="ts">
import Carousel from "./components/carousel.vue";
import StockTable from "./components/StockTable.vue";
</script>
