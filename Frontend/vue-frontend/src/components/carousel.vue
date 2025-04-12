<script setup lang="ts">
import { useCarouselStore } from "@/stores/carouselStore";
import { onMounted } from "vue";
import { useStockStore } from "@/stores/stockStore";

const carouselStore = useCarouselStore();
const stockStore = useStockStore();

onMounted(() => {
  carouselStore.carouselStore();
});

function seacrhClic(search) {
  stockStore.fetchStocks(search);
}
</script>

<template>
  <div
    class="slider"
    :style="{
      '--width': '200px',
      '--height': '100px',
      '--quantity': carouselStore.stockRecomendation.length,
    }"
  >
    <div v-if="carouselStore.loading" class="loading">
      Cargando recomendaciones...
    </div>
    <div v-else-if="carouselStore.error" class="error-message">
      {{ carouselStore.error }}
    </div>

    <div v-else class="list">
      <div
        v-for="(stock, index) in carouselStore.stockRecomendation"
        :key="stock.Id"
        class="item"
        :style="{
          '--position': index + 1,
          background: '#ffff',
        }"
      >
        <div class="card" @click="seacrhClic(stock.Ticker)">
          <h1 style="color: black">({{ stock.Ticker }})</h1>
          <div class="target-range">
            <span class="label">Target:</span>
            <span style="color: black"
              >${{ stock.TargetTo }} - ${{ stock.TargetTo }}</span
            >
          </div>
          <div class="potential-gain">
            <span class="label">Potencial:</span>
            <span
              :class="{
                positive: stock.PotentialGain >= 0,
                negative: stock.PotentialGain < 0,
              }"
            >
              {{ stock.PotentialGain }}%
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.slider {
  width: 100%;
  height: var(--height);
  overflow: hidden;
  mask-image: linear-gradient(to right, transparent, #000 10% 90%, transparent);
}

.list {
  display: flex;
  width: 100%;
  min-width: calc(var(--width) * var(--quantity));
  position: relative;
}

.item {
  width: var(--width);
  height: var(--height);
  position: absolute;
  left: 100%;
  animation: autoRun 10s linear infinite;
  transition: filter 0.5s;
  animation-delay: calc(
    (10s / var(--quantity)) * (var(--position) - 1) - 10s
  ) !important;
}

@keyframes autoRun {
  from {
    left: 100%;
  }
  to {
    left: calc(var(--width) * -1);
  }
}

.slider:hover .item {
  animation-play-state: paused !important;
  filter: grayscale(1);
}

.item:hover {
  filter: grayscale(0);
  transform: scale(1.02);
  z-index: 1;
}

.card {
  width: 100%;
  height: 100%;
  padding: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.5);
  color: white;
  text-align: center;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  backdrop-filter: blur(2px);
}

.card p {
  margin: 0;
  font-size: 13px;
  line-height: 1.4;
}

.company {
  font-weight: 600;
  font-size: 15px !important;
}

.brokerage {
  color: #c4c4c4;
  font-size: 12px !important;
}

.target-range,
.potential-gain {
  display: flex;
  justify-content: space-between;
  padding: 0 10px;
}

.label {
  color: #000000;
}

.positive {
  color: #4caf50;
}

.negative {
  color: #f44336;
}

.time {
  margin-top: 8px !important;
  font-size: 11px !important;
  color: #a0a0a0;
}

.loading,
.error-message {
  text-align: center;
  padding: 20px;
  color: white;
  font-size: 16px;
}

.error-message {
  color: #ff5555;
}
</style>
