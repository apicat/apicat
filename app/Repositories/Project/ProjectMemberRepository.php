<?php


namespace App\Repositories\Project;

use App\Repositories\User\UserRepository;
use Illuminate\Support\Carbon;
use App\Models\ProjectMember;
use App\Models\ProjectShare;

class ProjectMemberRepository
{
    /**
     * 检查成员是否在项目中
     * @param int $projectID 项目id
     * @param int $userID 成员id
     * @return boolean
     */
    public static function inThisProject(int $projectID, int $userID)
    {
        return (bool)ProjectMember::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->value('id');
    }

    /**
     * 获取成员项目权限
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return int
     */
    public static function getAuthority(int $projectID, int $userID)
    {
        return ProjectMember::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->value('authority');
    }

    /**
     * 添加项目成员
     * @param int $projectID 项目id
     * @param int $userID 团队成员id
     * @param int $authority 成员权限:0管理者,1维护者,2阅读者
     * @return mixed
     */
    public static function create($projectID, $userID, $authority = 2)
    {
        // 先查询该成员是否加入过项目
        $member = ProjectMember::withTrashed()->where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->first();
        if ($member) {
            if (!$member->trashed()) {
                // 成员当前未被移出项目
                return $member;
            }

            // 成员加入过项目，直接恢复
            $member->created_at = Carbon::now();
            $member->updated_at = Carbon::now();
            $member->authority = $authority;
            $member->restore();
            return $member;
        }

        // 成员之前未加入过该项目
        return ProjectMember::create([
            'project_id' => $projectID,
            'user_id' => $userID,
            'authority' => $authority
        ]);
    }

    /**
     * 成员是否有维护权限
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return boolean
     */
    public static function hasAuthority($projectID, $userID)
    {
        $authority = ProjectMember::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->value('authority');

        if ($authority === null) {
            return false;
        }

        return 2 > $authority;
    }

    /**
     * 检查成员是否是项目管理者
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return boolean 是: true  不是: false
     */
    public static function isAdmin($projectID, $userID)
    {
        return 0 === ProjectMember::where([
                ['project_id', $projectID],
                ['user_id', $userID]
            ])->value('authority');
    }

    /**
     * 修改成员权限
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @param int $authority 权限:0管理者,1维护者,2阅读者
     * @return boolean 成功: true  失败: false
     */
    public static function changeAuthority($projectID, $userID, $authority)
    {
        return (bool)ProjectMember::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->update(['authority' => $authority]);
    }

    /**
     * 从项目中删除成员
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return boolean
     */
    public static function remove($projectID, $userID)
    {
        ProjectShare::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->delete();

        return (bool)ProjectMember::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->delete();
    }

    /**
     * 项目成员数量
     * @param int $projectID 项目id
     * @return int
     */
    public static function projectMemberCount($projectID)
    {
        return ProjectMember::where('project_id', $projectID)->count();
    }

    /**
     * 项目成员列表
     * @param int $projectID 项目id
     * @param int $offset 游标起始位置
     * @param int $limit 获取数量
     * @return mixed
     */
    public static function projectMembers($projectID, $offset, $limit)
    {
        return ProjectMember::where('project_id', $projectID)->offset($offset)->limit($limit)->latest()->get();
    }

    /**
     * 项目中所有用户id
     *
     * @param int $projectID 项目id
     * @return array
     */
    public static function projectAllUserIds($projectID)
    {
        return ProjectMember::where('project_id', $projectID)->pluck('user_id')->toArray();
    }

    /**
     * 批量添加成员
     *
     * @param int $projectID 项目id
     * @param array $userIds 用户id数组
     * @param integer $authority 成员权限:0管理者,1维护者,2阅读者
     * @return boolean
     */
    public static function batchCreate($projectID, $userIds, $authority = 2)
    {
        if (count($userIds) == 1) {
            if (!is_numeric($userIds[0])) {
                return false;
            }

            if (!UserRepository::checkUserExists($userIds[0])) {
                return false;
            }

            return (bool)self::create($projectID, $userIds[0], $authority);
        }

        $allUserIds = UserRepository::userIds();

        $successNum = 0;

        foreach ($userIds as $userId) {
            if (!is_numeric($userId)) {
                continue;
            }

            if (!in_array($userId, $allUserIds)) {
                // 成员不存在
                continue;
            }

            if (self::create($projectID, $userId, $authority)) {
                $successNum++;
            }
        }

        return (bool)$successNum;
    }
}