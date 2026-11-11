import { request } from '../request';

/** 已上传图片 */
export interface UploadedImage {
  /** 相对路径（写库用） */
  path: string;
  /** 绝对地址（预览用） */
  url: string;
  name: string;
  size: number;
}

/**
 * 后台上传图片（site_user 令牌）
 *
 * 用于施工日志现场照片、变更附件等。
 *
 * @param file 图片文件
 */
export function fetchUploadAdminImage(file: File) {
  const formData = new FormData();
  formData.append('file', file);

  return request<{ files: UploadedImage[] }>({
    url: '/upload/adminImage',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  });
}

export function fetchUploadAdminContractFile(file: File) {
  const formData = new FormData();
  formData.append('file', file);

  return request<{ files: UploadedImage[] }>({
    url: '/upload/adminContractFile',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' }
  });
}
