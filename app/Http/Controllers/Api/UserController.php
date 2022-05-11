<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Modules\Image\Uploader;
use App\Repositories\User\UserRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class UserController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    /**
     * 个人信息
     *
     * @return array
     */
    public function profile()
    {
        $authorityName = ['超级管理员', '管理员', '普通成员'];

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'user_id' => Auth::id(),
                'name' => Auth::user()->name,
                'email' => Auth::user()->email,
                'avatar' => Auth::user()->avatar ? asset(Auth::user()->avatar)  : '',
                'authority' => Auth::user()->authority,
                'authority_name' => $authorityName[Auth::user()->authority]
            ]
        ];
    }

    /**
     * 更换头像
     *
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeAvatar(Request $request)
    {
        $request->validate([
            'avatar' => 'required|image|mimes:jpeg,jpg,png|max:1024',
            'cropped_x' => 'required|integer|min:0',
            'cropped_y' => 'required|integer|min:0',
            'cropped_width' => 'required|integer|min:1',
            'cropped_height' => 'required|integer|min:1',
        ]);

        if (!$request->hasFile('avatar')) {
            throw ValidationException::withMessages([
                'avatar' => '请上传头像图片',
            ]);
        }

        if (!$request->file('avatar')->isValid()) {
            throw ValidationException::withMessages([
                'avatar' => '无效的头像图片',
            ]);
        }

        if ($request->input('cropped_width') != $request->input('cropped_height')) {
            throw ValidationException::withMessages([
                'avatar' => '图片裁减信息有误',
            ]);
        }

        $filename = md5(Auth::id() . '_' . time()) . '.' . $request->file('avatar')->extension();
        $uploader = new Uploader($request->file('avatar')->path(), $filename, '/images/avatars/');

        // 检查图片是否被正确加载
        if (!$uploader->imgMakeResult()) {
            throw ValidationException::withMessages([
                'avatar' => '无效的头像图片',
            ]);
        }

        $res = $uploader->croppedSet(
            $request->input('cropped_width'),
            $request->input('cropped_height'),
            $request->input('cropped_x'),
            $request->input('cropped_y')
        );
        if (!$res) {
            throw ValidationException::withMessages([
                'avatar' => '图片裁减信息有误',
            ]);
        }

        $uploader->resizeSet(200, 200);

        $url = $uploader->save();

        if (UserRepository::editMyAccountInfo(['avatar' => $url])) {
            return [
                'status' => 0,
                'msg' => '头像上传成功',
                'data' => asset($url)
            ];
        }

        throw ValidationException::withMessages([
            'result' => '头像上传失败，请稍后重试。',
        ]);
    }

    /**
     * 修改个人信息
     *
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeProfile(Request $request)
    {
        $request->validate([
            'name' => ['required', 'string', 'max:30'],
            'email' => ['required', 'string', 'email', 'max:255'],
        ]);

        if (UserRepository::checkEmailExists($request->input('email'), Auth::id())) {
            throw ValidationException::withMessages([
                'result' => '邮箱已被使用，请更换邮箱后再试。',
            ]);
        }

        $data = [];

        if ($request->input('name') != Auth::user()->name) {
            $data = ['name' => $request->input('name')];
        }

        if ($request->input('email') != Auth::user()->email) {
            $data['email'] = $request->input('email');
        }

        if (!$data) {
            return ['status' => 0, 'msg' => '修改成功'];
        }

        if (UserRepository::editMyAccountInfo($data)) {
            return ['status' => 0, 'msg' => '修改成功'];
        }

        throw ValidationException::withMessages([
            'result' => '修改个人信息失败，请稍后重试。',
        ]);
    }

    /**
     * 修改密码
     *
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changePassword(Request $request)
    {
        $request->validate([
            'password' => 'required|min:8',
            'new_password' => 'required|confirmed|min:8',
            'new_password_confirmation' => 'required|min:8',
        ]);

        if (!UserRepository::checkMyPassword($request->input('password'))) {
            throw ValidationException::withMessages([
                'password' => '当前密码不正确',
            ]);
        }

        if ($request->input('password') == $request->input('new_password')) {
            throw ValidationException::withMessages([
                'password' => '当前密码与新密码一致',
            ]);
        }

        if (UserRepository::changeMyPassword($request->input('new_password'))) {
            return ['status' => 0, 'msg' => '密码修改成功'];
        }

        throw ValidationException::withMessages([
            'result' => '密码修改失败，请稍后重试。',
        ]);
    }
}
