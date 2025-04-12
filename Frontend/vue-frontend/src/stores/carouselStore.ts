import { defineStore } from "pinia";

export interface StockRecomendation {
  Id: number;
  Ticker: string;
  TargetFrom: string;
  TargetTo: string;
  Company: string;
  Action: string;
  Brokerage: string;
  RatingFrom: string;
  RatingTo: string;
  Time: string;
  PotentialGain: number;
}
const apiUrl = import.meta.env.VITE_API_URL;

export const useCarouselStore = defineStore("carouselStore", {
  state: () => ({
    stockRecomendation: [] as StockRecomendation[],
    loading: false,
    error: null as string | null,
    limit: 10,
  }),

  actions: {
    async carouselStore() {
      this.loading = true;
      try {
        const url = new URL(apiUrl + "/stock-recomendation");
        const response = await fetch(url.toString());
        if (!response.ok) {
          throw new Error(`Error en la solicitud: ${response.statusText}`);
        }
        const data = await response.json();
        this.stockRecomendation = data.stocks;
        console.log(this.stockRecomendation);
      } catch (error) {
        this.error =
          error instanceof Error ? error.message : "Error desconocido";
      } finally {
        this.loading = false;
      }
    },
  },
});
