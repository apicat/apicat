<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Project\ProjectMemberRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\User\UserRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class ProjectMemberController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 项目成员分页列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $request->validate([
            'page' => 'nullable|integer|min:1'
        ]);

        $page = $request->has('page') ? $request->input('page') : 1;
        $limit = 15;
        $offset = ($page - 1) * $limit;

        $authority = ProjectMemberRepository::getAuthority(ProjectRepository::active()->id, Auth::id());
        $isAdmin = $authority == 0;

        $projectMemberCount = ProjectMemberRepository::projectMemberCount(ProjectRepository::active()->id);
        if ($projectMemberCount < 1) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => [
                    'page' => 1,
                    'total_page' => 1,
                    'total_members' => 0,
                    'is_admin' => $isAdmin,
                    'project_members' => []
                ]
            ];
        }

        $totalPage = ceil($projectMemberCount / $limit);
        if ($totalPage < $page) {
            $page = $totalPage;
            $offset = ($page - 1) * $limit;
        }

        $projectMemberArr = [];
        $projectMembers = ProjectMemberRepository::projectMembers(ProjectRepository::active()->id, $offset, $limit);
        if ($projectMembers->count() > 0) {
            $userIds = [];
            foreach ($projectMembers as $member) {
                $userIds[] = $member->user_id;
            }

            $userArr = [];
            $users = UserRepository::getUsersByUserID($userIds);

            foreach ($users as $user) {
                $userArr[$user->id] = [
                    'id' => $user->id,
                    'avatar' => $user->avatar ?: '',
                    'name' => $user->name,
                    'email' => $user->email,
                ];
            }

            $authorityName = ['管理者', '维护者', '阅读者'];

            foreach ($projectMembers as $member) {
                $projectMemberArr[] = [
                    'user_id' => $member->user_id,
                    'avatar' => $userArr[$member->user_id]['avatar'],
                    'name' => $userArr[$member->user_id]['name'],
                    'email' => $userArr[$member->user_id]['email'],
                    'authority' => $member->authority,
                    'authority_name' => $authorityName[$member->authority]
                ];
            }
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'page' => $page,
                'total_page' => $totalPage,
                'total_members' => $projectMemberCount,
                'is_admin' => $isAdmin,
                'project_members' => $projectMemberArr
            ]
        ];
    }

    /**
     * 项目成员全部列表(id name avatar)
     * @return array
     */
    public function memberUserinfoList()
    {
        $userIds = ProjectMemberRepository::projectAllUserIds(ProjectRepository::active()->id);
        $myUserIDIndex = array_search(Auth::id(), $userIds);
        if ($myUserIDIndex !== false) {
            unset($userIds[$myUserIDIndex]);
        }

        if (empty($userIds)) {
            return ['status' => 0, 'msg' => '', 'data' => []];
        }

        $result = [];

        $users = UserRepository::getUsersByUserID($userIds);

        if ($users->count() > 0) {
            foreach ($users as $user) {
                $result[] = [
                    'user_id' => $user->id,
                    'name' => $user->name,
                    'avatar' => $user->avatar ? asset($user->avatar) : '',
                ];
            }
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $result
        ];
    }

    /**
     * 不在此项目的成员列表
     * @return array
     */
    public function notInProject()
    {
        $projectMemberIds = ProjectMemberRepository::projectAllUserIds(ProjectRepository::active()->id);

        $notInProjectUserArr = [];

        $users = UserRepository::users(0, 99999);

        if ($users->count() > 0) {
            foreach ($users as $user) {
                if (!in_array($user->id, $projectMemberIds)) {
                    $notInProjectUserArr[] = [
                        'user_id' => $user->id,
                        'name' => $user->name,
                        'email' => $user->email
                    ];
                }
            }
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $notInProjectUserArr
        ];
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
            'user_ids' => ['required', 'array', 'min:1'],
            'authority' => ['required', 'integer', 'in:1,2']
        ]);

        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有添加成员权限。',
            ]);
        }

        if (ProjectMemberRepository::batchCreate($request->input('project_id'), $request->input('user_ids'), $request->input('authority'))) {
            return [
                'status' => 0,
                'msg' => '成员添加成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '成员添加失败，请稍后重试。',
        ]);
    }

    /**
     * 修改成员权限
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeMemberAuthority(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1'],
            'authority' => ['required', 'integer', 'in:1,2']
        ]);

        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有权限修改成员权限。',
            ]);
        }

        if (ProjectMemberRepository::changeAuthority($request->input('project_id'), $request->input('user_id'), $request->input('authority'))) {
            return [
                'status' => 0,
                'msg' => '权限修改成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '权限修改失败，请稍后重试。',
        ]);
    }

    /**
     * 移出成员
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function removeMember(Request $request)
    {
        $request->validate([
            'user_id' => ['required', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->isAdmin()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有权限移出该成员。',
            ]);
        }

        if (ProjectMemberRepository::remove($request->input('project_id'), $request->input('user_id'))) {
            return [
                'status' => 0,
                'msg' => '移出成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '移出失败，请稍后重试。',
        ]);
    }
}
