import type { CreateAxiosDefaults } from 'axios';
import type { IAxiosRetryConfig } from 'axios-retry';
import { stringify } from 'qs';
import { isHttpSuccess } from './shared';
import type { RequestOption } from './type';

export function createDefaultOptions<
  ResponseData,
  ApiData = ResponseData,
  State extends Record<string, unknown> = Record<string, unknown>
>(options?: Partial<RequestOption<ResponseData, ApiData, State>>) {
  const opts: RequestOption<ResponseData, ApiData, State> = {
    defaultState: {} as State,
    transform: async response => response.data as unknown as ApiData,
    transformBackendResponse: async response => response.data as unknown as ApiData,
    onRequest: async config => config,
    isBackendSuccess: _response => true,
    onBackendFail: async () => {},
    onError: async () => {}
  };

  if (options?.transform) {
    opts.transform = options.transform;
  } else {
    opts.transform = options?.transformBackendResponse || opts.transform;
  }

  Object.assign(opts, options);

  return opts;
}

export function createRetryOptions(config?: Partial<CreateAxiosDefaults>) {
  const retryConfig: IAxiosRetryConfig = {
    retries: 0
  };

  Object.assign(retryConfig, config);

  return retryConfig;
}

export function createAxiosConfig(config?: Partial<CreateAxiosDefaults>) {
  const TEN_SECONDS = 10 * 1000;

  const axiosConfig: CreateAxiosDefaults = {
    timeout: TEN_SECONDS,
    headers: {
      'Content-Type': 'application/json'
    },
    validateStatus: isHttpSuccess,
    /**
     * Query 参数序列化
     *
     * 这里必须加 `skipNulls`：qs 默认会把 null 序列化成空串（`status=`），
     * 后端 GoFrame 会把空串绑定成「指向 0 的非空指针」，于是变成 `WHERE status = 0`，
     * 条件筛选类接口会一条都查不出来。skipNulls 让 null 参数直接不出现在 URL 上，
     * 语义上就是「该筛选项未传」，与业务预期一致。
     */
    paramsSerializer: params => {
      return stringify(params, { skipNulls: true });
    }
  };

  Object.assign(axiosConfig, config);

  return axiosConfig;
}
