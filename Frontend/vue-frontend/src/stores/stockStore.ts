import { defineStore } from "pinia";
const apiUrl = import.meta.env.VITE_API_URL;
export interface Stock {
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
}

export const useStockStore = defineStore("stockStore", {
  state: () => ({
    stocks: [] as Stock[],
    loading: false,
    error: null as string | null,
    limit: 10,
    totalStocks: 1,
    currentPageKey: 1,
    rowsPerPage: 10,
  }),

  actions: {
    async fetchStocks(
      searchKey?: string,
      limit: string = this.limit.toString()
    ) {
      this.loading = true;
      try {
        const url = new URL(apiUrl + "/stock");

        if (searchKey) {
          url.searchParams.append("search", searchKey);
        }

        url.searchParams.append("limit", limit);
        url.searchParams.append("currentPage", this.currentPageKey.toString());

        const response = await fetch(url.toString());
        if (!response.ok) {
          throw new Error(`Error en la solicitud: ${response.statusText}`);
        }

        const data = await response.json();
        this.stocks = data.stocks;
        this.totalStocks = data.totalCount;
      } catch (error) {
        this.error =
          error instanceof Error ? error.message : "Error desconocido";
      } finally {
        this.loading = false;
      }
    },

    setRowsPerPage(rows: number) {
      this.limit = rows;
    },
  },
});
