declare namespace Api {
  /**
   * namespace Print
   *
   * backend api module: "print"（打印设置）
   */
  namespace Print {
    /** 打印设置 */
    interface PrintSetting {
      /** 抬头标题，空则打印时仅显示单据名 */
      title: string;
      /** 副标题 */
      subTitle: string;
      /** 页脚文字 */
      footer: string;
      /** 是否显示打印日期 */
      showDate: boolean;
      /** 是否加盖电子章 */
      showSeal: boolean;
      /** 印章图片相对路径 */
      sealImage: string;
      /** 印章图片绝对地址（仅查询返回） */
      sealImageUrl?: string;
    }

    /** 打印设置保存参数 */
    interface PrintSettingForm {
      title: string;
      subTitle: string;
      footer: string;
      showDate: boolean;
      showSeal: boolean;
      sealImage: string;
    }
  }
}
