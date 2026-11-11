import { ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import dayjs from 'dayjs';

const STORAGE_KEY = 'report-year';

function normalizeYear(value: unknown) {
  const year = Number(value);

  return year >= 2000 && year <= 2100 ? year : 0;
}

/** 报表年份：优先读 URL query，其次沿用上次选择，最后取当前年份 */
export function useReportYear() {
  const route = useRoute();
  const router = useRouter();

  const storedYear = normalizeYear(sessionStorage.getItem(STORAGE_KEY));
  const year = ref(normalizeYear(route.query.year) || storedYear || dayjs().year());

  watch(year, value => {
    sessionStorage.setItem(STORAGE_KEY, String(value));
    void router.replace({ query: { ...route.query, year: String(value) } });
  });

  return { year };
}
