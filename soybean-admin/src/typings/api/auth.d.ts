declare namespace Api {
  /**
   * namespace Auth
   *
   * backend api module: "auth"
   */
  namespace Auth {
    interface LoginToken {
      token: string;
      refreshToken: string;
    }

    interface UserInfo {
      userId: number;
      userName: string;
      /** 昵称/姓名 */
      nickname: string;
      /** 头像URL */
      avatar: string;
      /** 手机号 */
      phone: string;
      /** 角色：1=管理员 2=员工 3=商家 */
      role: number;
      /** 所属商家ID */
      merchantId: number;
      /** 所属公司ID */
      companyId: number;
      /** Whether finance setup is complete */
      finInitialized: boolean;
      roles: string[];
      buttons: string[];
    }
  }
}
