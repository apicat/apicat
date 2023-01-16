<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Iteration\StarRepository;
use App\Repositories\Project\GroupRepository;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\ProjectShareRepository;
use App\Repositories\User\UserRepository;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Illuminate\Support\Facades\Auth;

class ProjectController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 开启关闭项目分享
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function share(Request $request)
    {
        $request->validate([
            'share' => ['required', 'integer', 'in:0,1']
        ]);

        if (!ProjectMemberRepository::hasAuthority($request->input('project_id'), Auth::id())) {
            throw ValidationException::withMessages([
                'project_id' => '无法分享该项目',
            ]);
        }

        if ($request->input('share')) {
            if (!$share = ProjectShareRepository::startShare($request->input('project_id'), Auth::id())) {
                throw ValidationException::withMessages([
                    'result' => '分享项目失败，请稍后重试。',
                ]);
            }

            return [
                'status' => 0,
                'msg' => '开启分享成功',
                'data' => [
                    'link' => route('app.index', ['projectID' => ProjectRepository::active()->id]),
                    'secret_key' => $share->secret_key
                ]
            ];
        } else {
            if (!ProjectShareRepository::closeShare($request->input('project_id'), Auth::id())) {
                throw ValidationException::withMessages([
                    'result' => '关闭项目分享失败，请稍后重试。',
                ]);
            }

            return [
                'status' => 0,
                'msg' => '关闭项目分享成功'
            ];
        }
    }

    /**
     * 重置分享项目访问秘钥
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function resetShareSecretKey(Request $request)
    {
        if (!$secretKey = ProjectShareRepository::changeSharePassword($request->input('project_id'), Auth::id())) {
            throw ValidationException::withMessages([
                'result' => '修改失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '修改成功',
            'data' => $secretKey
        ];
    }

    /**
     * 项目设置
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function storeSetting(Request $request)
    {
        $request->validate([
            'icon_link' => ['nullable', 'string', 'max:255'],
            'name' => ['required', 'string', 'min:2', 'max:255'],
            'visibility' => ['required', 'int', 'in:0,1'],
            'description' => ['nullable', 'string', 'max:255'],
        ]);

        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有修改权限。',
            ]);
        }

        if (ProjectRepository::edit($request->input('project_id'), $request->all())) {
            return [
                'status' => 0,
                'msg' => '项目修改成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '项目修改失败，请稍后重试。',
        ]);
    }

    /**
     * 项目分组
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeGroup(Request $request)
    {
        $request->validate([
            'old_group_id' => 'required|integer|min:0',
            'new_group_id' => 'required|integer|min:0',
        ]);

        if ($request->input('new_group_id') > 0) {
            $newGroup = GroupRepository::find($request->input('new_group_id'));
            if (!$newGroup or $newGroup->user_id != Auth::id()) {
                throw ValidationException::withMessages([
                    'new_group_id' => '分组不存在',
                ]);
            }
        }

        if ($request->input('old_group_id') > 0) {
            $oldGroup = GroupRepository::find($request->input('old_group_id'));
            if (!$oldGroup or $oldGroup->user_id != Auth::id()) {
                throw ValidationException::withMessages([
                    'new_group_id' => '分组信息有误',
                ]);
            }
        }

        if (!GroupRepository::changeGroup($request->input('project_id'), $request->input('old_group_id'), $request->input('new_group_id'))) {
            throw ValidationException::withMessages([
                'result' => '移动分组失败，请稍后重试。',
            ]);
        }

        if ($request->input('new_group_id') > 0) {
            return ['status' => 0, 'msg' => '项目已移动至' . $newGroup->name . '分组'];
        }

        return ['status' => 0, 'msg' => '项目已移出分组'];
    }

    /**
     * 修改项目拥有者
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeOwner(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1']
        ]);

        if (!UserRepository::isSuperAdmin(Auth::user())) {
            throw ValidationException::withMessages([
                'result' => '您不是团队超级管理员，无权修改项目拥有者。',
            ]);
        }

        if (!UserRepository::checkUserExists($request->input('user_id'))) {
            throw ValidationException::withMessages([
                'user_id' => '用户不存在',
            ]);
        }

        if (ProjectRepository::active()->user_id == $request->input('user_id')) {
            return [
                'status' => 0,
                'msg' => '修改成功'
            ];
        }

        if (!ProjectRepository::changeOwner($request->input('project_id'), $request->input('user_id'), ProjectRepository::active()->user_id)) {
            throw ValidationException::withMessages([
                'result' => '修改失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '修改成功'
        ];
    }

    /**
     * 移交项目
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function transfer(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '您不是项目管理者，无权移交项目。',
            ]);
        }

        // 项目只能移交给维护者
        if (!ProjectMemberRepository::hasAuthority($request->input('project_id'), $request->input('user_id'))) {
            throw ValidationException::withMessages([
                'result' => '该成员不是项目维护者，无权接收项目。',
            ]);
        }

        if (ProjectRepository::changeOwner($request->input('project_id'), $request->input('user_id'), Auth::id())) {
            return [
                'status' => 0,
                'msg' => '移交项目成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '移交项目失败，请稍后重试。',
        ]);
    }

    /**
     * 退出项目
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function quit(Request $request)
    {
        if (ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'result' => '由于您是项目管理者，所以您无法退出项目，请先将项目移交给其他成员后即可退出。',
            ]);
        }

        if (ProjectMemberRepository::remove($request->input('project_id'), Auth::id())) {
            StarRepository::del($request->input('project_id'));

            return [
                'status' => 0,
                'msg' => '成功退出项目'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '退出项目失败，请稍后重试。',
        ]);
    }

    /**
     * 删除项目
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function remove(Request $request)
    {
        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有权限删除该项目。',
            ]);
        }

        if (ProjectRepository::remove($request->input('project_id'))) {
            StarRepository::delAll($request->input('project_id'));

            return [
                'status' => 0,
                'msg' => '项目删除成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '项目删除失败，请稍后重试。',
        ]);
    }
}
