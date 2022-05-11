<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Foundation\Auth\ThrottlesLogins;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;
use Illuminate\Validation\ValidationException;
use App\Models\User;
use App\Modules\OperateThrottles\EmailRegisterThrottle;

/**
 * 邮箱登录、注册处理
 *
 * @package App\Http\Controllers\Auth
 */
class EmailController extends Controller
{
    use ThrottlesLogins;

    /**
     * api访问令牌
     *
     * @var string
     */
    protected string $apiToken;

    /**
     * 注册
     *
     * @param Request $request
     * @return JsonResponse|Response
     */
    public function register(Request $request)
    {
        $request->validate([
            'email' => 'required|string|email|max:255|unique:users',
            'password' => 'required|string|min:8'
        ]);

        $throttle = new EmailRegisterThrottle;

        if ($throttle->hasTooManyAttempts($request)) {
            $throttle->fireLockoutEvent($request);

            $seconds = $throttle->limiter()->availableIn(
                $throttle->throttleKey($request)
            );

            throw ValidationException::withMessages([
                'result' => '您尝试的注册次数过多，请 ' . $seconds . ' 秒后再试。'
            ]);
        }

        $user = User::create([
            'name' => explode('@', $request->input('email'))[0],
            'email' => $request->input('email'),
            'password' => Hash::make($request->input('password')),
            'authority' => 2
        ]);
        if (!$user) {
            return response()->json(['status' => -1, 'msg' => '注册失败，请稍后重试。']);
        }

        if ($user->id === 1) {
            $user->authority = 0;
            $user->save();
        }

        $throttle->incrementAttempts($request);

        if (!$apiToken = Auth::guard('api')->login($user)) {
            return response()->json(['status' => -1, 'msg' => '注册失败，请稍后重试。']);
        }

        return response()->json(['status' => 0, 'msg' => '注册成功'])->header('Authorization', $apiToken);
    }

    /**
     * 登录
     *
     * @param Request $request
     * @return Response|void
     * @throws ValidationException
     */
    public function login(Request $request)
    {
        $request->validate([
            'email' => 'required|string|email|max:255',
            'password' => 'required|string|min:8'
        ]);

        if (method_exists($this, 'hasTooManyLoginAttempts') &&
            $this->hasTooManyLoginAttempts($request)) {
            $this->fireLockoutEvent($request);

            return $this->sendLockoutResponse($request);
        }

        if ($this->attemptLogin($request)) {
            return $this->sendLoginResponse($request);
        }

        $this->incrementLoginAttempts($request);

        return $this->sendFailedLoginResponse();
    }

    /**
     * 尝试登录
     *
     * @param Request $request
     * @return boolean
     */
    protected function attemptLogin(Request $request)
    {
        if (!$user = User::where('email', $request->input('email'))->first()) {
            return false;
        }

        if (!Hash::check($request->input('password'), $user->password)) {
            return false;
        }

        if (!$this->apiToken = Auth::guard('api')->login($user)) {
            return false;
        }

        return true;
    }

    /**
     * 返回登录成功信息
     *
     * @param Request $request
     * @return JsonResponse
     */
    protected function sendLoginResponse(Request $request)
    {
        $this->clearLoginAttempts($request);
        return response()->json(['status' => 0, 'msg' => '登录成功'])->header('Authorization', $this->apiToken);
    }

    /**
     * 返回登录失败信息
     *
     * @throws ValidationException
     */
    protected function sendFailedLoginResponse()
    {
        throw ValidationException::withMessages([
            'email' => [trans('auth.email')],
        ]);
    }

    /**
     * 获取登录的用户名
     *
     * @return void
     */
    protected function username()
    {
        return 'email';
    }
}