<?php

namespace App\Repositories\Project;

use App\Models\Project;
use App\Models\ProjectMember;
use App\Models\ProjectShare;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\DB;

class ProjectRepository
{
    public static $active;

    /**
     * 激活项目为当前操作项目，获取当前项目信息
     * @param Project|null $project
     * @return ActiveProjectRepository|false
     */
    public static function active($project = null)
    {
        if (!self::$active) {
            if (!$project) {
                return false;
            }
            self::$active = new ActiveProjectRepository($project);
        }
        return self::$active;
    }

    /**
     * 用户参与的所有项目列表
     * @param int $userID 用户id
     * @return Collection
     */
    public static function list($userID)
    {
        $projectIds = [];
        $projects = ProjectMember::where('user_id', $userID)->get();
        if ($projects->count() < 1) {
            // 成员没有参与任何项目
            return new Collection;
        }

        $authorityNames = ['管理者', '维护者', '阅读者'];
        $authorities = [];
        foreach ($projects as $project) {
            $projectIds[] = $project->project_id;
            $authorities[$project->project_id] = [
                'authority' => $project->authority,
                'authority_name' => $authorityNames[$project->authority]
            ];
        }

        $shares = ProjectShare::whereIn('project_id', $projectIds)->where('user_id', $userID)->get();
        $secretkeys = [];
        if ($shares->count() > 0) {
            foreach ($shares as $share) {
                $secretkeys[$share->project_id] = $share->secret_key;
            }
        }

        return self::recordByIds($projectIds, $authorities, $secretkeys);
    }

    /**
     * 通过项目id数组获取多条项目记录
     * @param array $ids 项目id数组
     * @param array $authorities 一个成员在对应项目中的权限信息
     * @param array $secretkeys 一个成员分享对应项目的访问秘钥
     * @return \Illuminate\Database\Eloquent\Collection
     */
    public static function recordByIds($ids, $authorities = [], $secretkeys = [])
    {
        $projects = Project::whereIn('id', $ids)->latest()->get();
        if ($projects->count() > 0) {
            foreach ($projects as $project) {
                $project->preview_link = route('app.index', ['projectID' => $project->id]);
                $project->default_link = $project->authority < 2 ? route('editor.doc', ['projectID' => $project->id]) : $project->preview_link;

                if ($authorities and isset($authorities[$project->id])) {
                    $project->authority = $authorities[$project->id]['authority'];
                    $project->authority_name = $authorities[$project->id]['authority_name'];
                }

                if ($secretkeys and isset($secretkeys[$project->id])) {
                    $project->secret_key = $secretkeys[$project->id];
                }
            }
        }

        return $projects;
    }

    /**
     * 获取项目信息
     * @param int $projectID 项目id
     * @return Project
     */
    public static function get($projectID)
    {
        return Project::find($projectID);
    }

    /**
     * 分组下的所有项目列表
     *
     * @param int $userID 成员id
     * @param int $groupID 分组id
     * @return \Illuminate\Database\Eloquent\Collection|Collection
     */
    public static function groupList($userID, $groupID)
    {
        $projectIds = GroupRepository::projectIds($groupID);
        if (!$projectIds) {
            // 分组里没有项目
            return new Collection;
        }

        $authorityNames = ['管理者', '维护者', '阅读者'];
        $authorities = [];

        // 没有被删除的项目id都放进此数组
        $existProjectIds = [];

        $members = ProjectMember::where('user_id', $userID)->whereIn('project_id', $projectIds)->get();
        if ($members->count() > 0) {
            foreach ($members as $member) {
                $existProjectIds[] = $member->project_id;

                $authorities[$member->project_id] = [
                    'authority' => $member->authority,
                    'authority_name' => $authorityNames[$member->authority]
                ];
            }
        } else {
            foreach ($projectIds as $projectID) {
                GroupRepository::changeGroup($projectID, $groupID, 0);
            }
        }

        if (!$authorities) {
            // 该成员不在分组里的项目中
            return new Collection;
        }

        if ($diff = array_diff($projectIds, $existProjectIds)) {
            // 如果 $projectIds 中的项目id在 $existProjectIds 中没有出现，说明对应id的项目已经被删除，
            // 但是在这个项目id作为脏数据还存在用户的分组里，这里可以帮助用户把脏数据清理掉。
            foreach ($diff as $projectID) {
                GroupRepository::changeGroup($projectID, $groupID, 0);
            }
        }

        $shares = ProjectShare::whereIn('project_id', $existProjectIds)->where('user_id', $userID)->get();
        $secretkeys = [];
        if ($shares->count() > 0) {
            foreach ($shares as $share) {
                $secretkeys[$share->project_id] = $share->secret_key;
            }
        }

        return self::recordByIds($existProjectIds, $authorities, $secretkeys);
    }

    /**
     * 当前所在团队中自己创建的项目列表
     *
     * @param int $userID 成员id
     * @return \Illuminate\Database\Eloquent\Collection|void 有项目返回项目集合，没有项目返回null
     */
    public static function ownList($userID)
    {
        $projectIds = ProjectMember::where([
            ['user_id', $userID],
            ['authority', 0]
        ])->pluck('project_id')->toArray();
        if ($projectIds) {
            return Project::whereIn('id', $projectIds)->get();
        }
    }

    /**
     * 当前所在团队中可编辑的项目列表
     *
     * @param int $userID 成员id
     * @return \Illuminate\Database\Eloquent\Collection|void 有项目返回项目集合，没有项目返回null
     */
    public static function writeList($userID)
    {
        $projectIds = ProjectMember::where([
            ['user_id', $userID],
            ['authority', '<', 2]
        ])->pluck('project_id')->toArray();
        if ($projectIds) {
            return Project::whereIn('id', $projectIds)->get();
        }
    }

    /**
     * 当前所在团队中只读的项目列表
     *
     * @param int $userID 成员id
     * @return \Illuminate\Database\Eloquent\Collection|void 有项目返回项目集合，没有项目返回null
     */
    public static function readList($userID)
    {
        $projectIds = ProjectMember::where([
            ['user_id', $userID],
            ['authority', 2]
        ])->pluck('project_id')->toArray();
        if ($projectIds) {
            return Project::whereIn('id', $projectIds)->get();
        }
    }

    /**
     * 创建项目
     * @param array $data 项目信息
     * @return Project|boolean  成功返回项目实例，失败返回false
     */
    public static function create($data)
    {
        $project = Project::create([
            'user_id' => Auth::id(),
            'icon' => $data['icon_link'] ?? '',
            'name' => $data['name'],
            'visibility' => $data['visibility'],
            'description' => $data['description'] ?? '',
        ]);
        if (!$project) {
            return false;
        }

        // 将创始人添加进项目成员表中
        $projectMember = ProjectMemberRepository::create($project->id, Auth::id(), 0);
        if (!$projectMember) {
            $project->delete();
            return false;
        }

        // 项目进行分组
        if (isset($data['group_id']) and $data['group_id'] > 0) {
            GroupRepository::changeGroup($project->id, 0, $data['group_id']);
        }

        return $project;
    }

    /**
     * 修改项目设置
     * @param int $projectID 项目id
     * @param array $data 修改信息
     * @return boolean 成功: true  失败: false
     */
    public static function edit($projectID, $data)
    {
        $saveData = [
            'name' => $data['name'],
            'visibility' => $data['visibility']
        ];

        if (isset($data['icon_link'])) {
            $saveData['icon'] = $data['icon_link'];
        }
        if (isset($data['description'])) {
            $saveData['description'] = $data['description'];
        }

        return (bool)Project::where('id', $projectID)->update($saveData);
    }

    /**
     * 修改项目拥有者
     * @param int $projectID 项目id
     * @param int $toUserID 新拥有者成员id
     * @param int $originUserID 原拥有者成员id
     * @return boolean 成功: true  失败: false
     */
    public static function changeOwner($projectID, $toUserID, $originUserID = 0)
    {
        if (!$originUserID) {
            $originUserID = ProjectMember::where([
                ['project_id', $projectID],
                ['authority', 0]
            ])->value('user_id');
        }

        if ($toUserID == $originUserID) {
            return true;
        }

        $inThisProject = ProjectMemberRepository::inThisProject($projectID, $toUserID);

        try {
            DB::transaction(function () use ($projectID, $toUserID, $originUserID, $inThisProject) {
                // 将项目新拥有者修改为新的管理者
                Project::where('id', $projectID)->update(['user_id' => $toUserID]);

                if (!$inThisProject) {
                    // 要转交的成员不在项目中，进行添加
                    ProjectMemberRepository::create($projectID, $toUserID, 0);
                } else {
                    // 要转交的成员在项目中，修改新管理者权限
                    ProjectMemberRepository::changeAuthority($projectID, $toUserID, 0);
                }

                // 修改原拥有者角色为阅读者
                ProjectMemberRepository::changeAuthority($projectID, $originUserID, 2);
            });
        } catch (\Exception $e) {
            return false;
        }
        return true;
    }

    /**
     * 删除项目
     * @param int $projectID 项目id
     * @return boolean 成功: true  失败: false
     */
    public static function remove($projectID)
    {
        if (!Project::where('id', $projectID)->delete()) {
            return false;
        }

        ProjectMember::where('project_id', $projectID)->delete();
        ProjectShare::where('project_id', $projectID)->delete();

        return true;
    }
}
