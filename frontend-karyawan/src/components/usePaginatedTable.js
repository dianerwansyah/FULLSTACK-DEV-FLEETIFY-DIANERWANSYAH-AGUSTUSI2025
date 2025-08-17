import { ref, watch } from "vue";
import api from "@/axios";

/**
 * Paginated table composable
 * @param {Object} options
 * @param {string} options.endpoint - API endpoint
 * @param {Function} options.normalizeFn - Normalizer for each item
 * @param {Object|Ref|ComputedRef} options.filters - Filter object or ref/computed
 */
export function usePaginatedTable({ endpoint, normalizeFn, filters = {} }) {
  const items = ref([]);
  const loading = ref(false);
  const totalItems = ref(0);

  const page = ref(1);
  const itemsPerPage = ref(10);

  const sortBy = ref([]);
  const sortDesc = ref(false);

  const resolveFilterValue = (val) => {
    if (val && typeof val === "object" && "value" in val) {
      return val.value;
    }
    return val;
  };

  const buildPayload = () => {
    const rawFilters =
      typeof filters === "function"
        ? filters()
        : typeof filters === "object" && "value" in filters
        ? filters.value
        : filters;

    const resolvedFilters = Object.fromEntries(
      Object.entries(rawFilters).map(([key, val]) => [
        key,
        resolveFilterValue(val),
      ])
    );

    return {
      page: page.value,
      per_page: itemsPerPage.value,
      sort_by: Array.isArray(sortBy.value) ? sortBy.value : [sortBy.value],
      filter: resolvedFilters,
    };
  };

  const fetchData = async () => {
    loading.value = true;
    try {
      const payload = buildPayload();
      const res = await api.post(endpoint, payload);

      const raw = res.data?.data || res.data?.items || [];
      items.value = Array.isArray(raw) ? raw.map(normalizeFn) : [];
      totalItems.value = res.data?.meta?.total || raw.length;
    } catch (err) {
      console.error("Fetch error:", err);
      items.value = [];
      totalItems.value = 0;
    } finally {
      loading.value = false;
    }
  };

  // Watch all reactive dependencies
  watch(
    [page, itemsPerPage, sortBy, sortDesc, ...Object.values(filters)],
    () => {
      fetchData();
    }
  );

  return {
    items,
    loading,
    totalItems,
    page,
    itemsPerPage,
    sortBy,
    sortDesc,
    fetchData,
  };
}
