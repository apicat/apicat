import {
  COMPLETE_INFO_NAME,
  COMPLETE_INFO_PATH,
  FORGETPASS_NAME,
  FORGETPASS_PATH,
  LOGIN_NAME,
  LOGIN_PATH,
  REGISTER_NAME,
  REGISTER_PATH,
  REGISTER_VERIFICATION_EMAIL_NAME,
  REGISTER_VERIFICATION_EMAIL_PATH,
  RESET_PASS_PATH,
  USER_CHANGE_ACCOUNT_EMAIL_NAME,
  USER_CHANGE_ACCOUNT_EMAIL_PATH,
} from './constant'

export const loginRoute = {
  path: LOGIN_PATH,
  name: LOGIN_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.login' },
  component: () => import('@/views/sign/LoginPage.vue'),
}

export const registerRoute = {
  path: REGISTER_PATH,
  name: REGISTER_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.register' },
  component: () => import('@/views/sign/RegisterPage.vue'),
}

export const forgetPassRoute = {
  path: FORGETPASS_PATH,
  name: FORGETPASS_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.forgetPass' },
  component: () => import('@/views/sign/ForgetPass.vue'),
}

export const resetPassRoute = {
  path: RESET_PASS_PATH,
  name: 'resetPass',
  meta: { ignoreAuth: true, title: 'app.pageTitles.resetPass' },
  component: () => import('@/views/sign/ResetPass.vue'),
}

export const registerVerificationEmailRoute = {
  path: REGISTER_VERIFICATION_EMAIL_PATH,
  name: REGISTER_VERIFICATION_EMAIL_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.verificationEmail' },
  component: () => import('@/views/sign/VerificationEmailForRegister.vue'),
}

export const verificationEmailForModifyRoute = {
  path: USER_CHANGE_ACCOUNT_EMAIL_PATH,
  name: USER_CHANGE_ACCOUNT_EMAIL_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.verificationEmail' },
  component: () => import('@/views/sign/VerificationEmailForModify.vue'),
}

export const completeInfoRoute = {
  path: COMPLETE_INFO_PATH,
  name: COMPLETE_INFO_NAME,
  meta: { ignoreAuth: true, title: 'app.pageTitles.completeInfo' },
  component: () => import('@/views/sign/CompleteInfo.vue'),
}
