import { request } from '../request';

/** 获取打印设置（未配置时返回默认值） */
export function fetchGetPrintSetting() {
  return request<{ setting: Api.Print.PrintSetting }>({
    url: '/print/setting/get',
    method: 'get'
  });
}

/** 保存打印设置 */
export function fetchSavePrintSetting(data: Api.Print.PrintSettingForm) {
  return request<null>({
    url: '/print/setting/save',
    method: 'post',
    data
  });
}
