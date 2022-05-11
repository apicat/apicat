<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\User\UserRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class TeamController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    /**
     * 成员列表
     *
     * @return array
     */
    public function members(Request $request)
    {
        $request->validate([
            'page' => 'nullable|integer|min:1',
            'page_size' => 'nullable|integer|min:1'
        ]);

        $page = $request->has('page') ? $request->input('page') : 1;
        $limit = $request->has('page_size') ? $request->input('page_size') : 15;
        $offset = ($page - 1) * $limit;
        $userCount = UserRepository::userCount();
        if ($userCount < 1) {
            // 没有团队成员
            return [
                'status' => 0,
                'msg' => '',
                'data' => [
                    'page' => 1,
                    'total_page' => 1,
                    'total_members' => 0,
                    'members' => []
                ]
            ];
        }

        $totalPage = ceil($userCount / $limit);
        if ($totalPage < $page) {
            $page = $totalPage;
            $offset = ($page - 1) * $limit;
        }

        $authorityNames = ['超级管理员', '管理员', '普通成员'];
        $memberArr = [];
        $users = UserRepository::users($offset, $limit);
        foreach ($users as $user) {
            $memberArr[] = [
                'user_id' => $user->id,
                'authority' => $user->authority,
                'authority_name' => $authorityNames[$user->authority],
                'name' => $user->name,
                'email' => $user->email,
            ];
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'page' => $page,
                'total_page' => $totalPage,
                'total_members' => $userCount,
                'members' => $memberArr
            ]
        ];
    }

    /**
     * 成员详情
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function member(Request $request)
    {
        $request->validate([
            'user_id' => 'required|integer|min:1'
        ]);

        $userInfo = UserRepository::getUserInfo($request->input('user_id'));
        if (!$userInfo) {
            throw ValidationException::withMessages([
                'user_id' => '成员不存在',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'user_id' => $userInfo->id,
                'authority' => $userInfo->authority,
                'name' => $userInfo->name,
                'avatar' => $userInfo->avatar ? asset($userInfo->avatar) : '',
                'email' => $userInfo->email
            ]
        ];
    }

    /**
     * 成员参与的项目列表
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function memberProjects(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1',]
        ]);

        if (!UserRepository::checkUserExists($request->input('user_id'))) {
            throw ValidationException::withMessages([
                'user_id' => '成员不存在',
            ]);
        }

        $projectArr = [];

        if (UserRepository::isSuperAdmin(Auth::id()) or (UserRepository::hasAuthority(Auth::id()) and !UserRepository::hasAuthority($request->input('user_id')))) {
            // 超级管理员可以查看所有人参与的项目
            // 管理员只能查看普通成员参与的项目
            $projects = ProjectRepository::list($request->input('user_id'));
            if ($projects->count() > 0) {
                foreach ($projects as $project) {
                    $projectArr[] = [
                        'id' => $project->id,
                        'preview_link' => $project->preview_link,
                        'default_link' => $project->default_link,
                        'icon' => $project->icon,
                        'name' => $project->name,
                        'visibility' => $project->visibility,
                        'authority' => $project->authority,
                        'authority_name' => $project->authority_name
                    ];
                }
            }
        }

        return ['status' => 0, 'msg' => '', 'data' => $projectArr];

    }

    /**
     * 添加成员
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function addMember(Request $request)
    {
        $request->validate([
            'email' => ['required', 'string', 'email', 'max:255', 'unique:users'],
            'name' => ['required', 'string', 'max:255'],
            'password' => ['required', 'string', 'min:8'],
            'authority' => ['required', 'integer', 'in:1,2']
        ]);

        if (!UserRepository::hasAuthority(Auth::user())) {
            // 只有团队管理员有修改权限
            throw ValidationException::withMessages([
                'result' => '您没有添加成员权限',
            ]);
        }

        if (UserRepository::isAdmin(Auth::user()) and $request->input('authority') == 1) {
            // 管理员只能修改普通成员的密码
            throw ValidationException::withMessages([
                'user_id' => '您没有添加管理成员权限',
            ]);
        }

        if (UserRepository::addAccount($request->all())) {
            return [
                'status' => 0,
                'msg' => '添加成员成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '添加成员失败，请稍后重试。',
        ]);
    }

    /**
     * 修改成员信息
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function editMemberInfo(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1'],
            'email' => ['required', 'email', 'max:255'],
            'name' => ['required', 'string', 'max:255'],
            'authority' => ['nullable', 'integer', 'in:1,2']
        ]);

        if (!UserRepository::hasAuthority(Auth::user())) {
            // 只有团队管理员有修改权限
            throw ValidationException::withMessages([
                'result' => '您没有修改权限',
            ]);
        }

        $data = [];

        if (UserRepository::isSuperAdmin(Auth::user())) {
            // 超级管理员
            $data['email'] = $request->input('email');
            $data['name'] = $request->input('name');
            $data['authority'] = $request->input('authority');
        } else {
            // 管理员
            if (UserRepository::hasAuthority((integer)$request->input('user_id'))) {
                // 管理员只能修改普通成员的信息
                throw ValidationException::withMessages([
                    'user_id' => '您没有修改权限',
                ]);
            }

            $data['email'] = $request->input('email');
            $data['name'] = $request->input('name');
        }

        if (UserRepository::checkEmailExists($request->input('email'), $request->input('user_id'))) {
            throw ValidationException::withMessages([
                'result' => '邮箱已被其他用户注册，请更换邮箱后再试。',
            ]);
        }

        if (UserRepository::editUserAccountInfo($request->input('user_id'), $data)) {
            return [
                'status' => 0,
                'msg' => '修改成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '修改成员信息失败，请稍后重试。',
        ]);
    }

    /**
     * 修改成员密码
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function editMemberPassword(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1'],
            'password' => ['required', 'string', 'min:8'],
        ]);

        if (!UserRepository::hasAuthority(Auth::user())) {
            // 只有团队管理员有修改权限
            throw ValidationException::withMessages([
                'result' => '您没有修改权限',
            ]);
        }

        if (UserRepository::isAdmin(Auth::user()) and UserRepository::hasAuthority((integer)$request->input('user_id'))) {
            // 管理员只能修改普通成员的密码
            throw ValidationException::withMessages([
                'user_id' => '您没有修改权限',
            ]);
        }

        if (UserRepository::editUserAccountInfo($request->input('user_id'), $request->all())) {
            return [
                'status' => 0,
                'msg' => '修改成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '修改成员密码失败，请稍后重试。',
        ]);
    }

    /**
     * 删除成员
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function removeMember(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1'],
        ]);

        // 管理权限检查
        if (!UserRepository::hasAuthority(Auth::user())) {
            throw ValidationException::withMessages([
                'result' => '您没有删除权限',
            ]);
        }

        if (UserRepository::isAdmin(Auth::user()) and UserRepository::hasAuthority((integer)$request->input('user_id'))) {
            // 管理员只能修改普通成员的密码
            throw ValidationException::withMessages([
                'user_id' => '您没有权限删除该成员',
            ]);
        }

        if (UserRepository::remove($request->input('user_id'))) {
            return [
                'status' => 0,
                'msg' => '删除成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '删除失败，请稍后重试。',
        ]);
    }
}
