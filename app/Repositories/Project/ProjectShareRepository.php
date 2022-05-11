<?php

namespace App\Repositories\Project;

use Illuminate\Support\Str;
use App\Models\ProjectShare;

class ProjectShareRepository
{
    /**
     * 检查项目是否被分享了
     * @param int $projectID 项目id
     * @return boolean true: 已被分享  false: 未被分享
     */
    public static function hasShared($projectID)
    {
        return (bool)ProjectShare::where('project_id', $projectID)->limit(1)->value('id');
    }

    /**
     * 通过成员id获取分享记录
     * @param int $projectID 项目id
     * @param int $userID 成员id
     * @return ProjectShare
     */
    public static function getByMemberID($projectID, $userID)
    {
        return ProjectShare::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->first();
    }

    /**
     * 校验项目访问密码
     * @param int $projectID 项目id
     * @param string $secret 访问密码
     * @return boolean
     */
    public static function check($projectID, $secret)
    {
        return (bool)ProjectShare::where([
            ['project_id', $projectID],
            ['secret_key', $secret]
        ])->value('id');
    }

    /**
     * 生成私有项目的分享
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return ProjectShare
     */
    public static function startShare($projectID, $userID)
    {
        $record = ProjectShare::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->first();

        if (!$record) {
            return ProjectShare::create([
                'project_id' => $projectID,
                'user_id' => $userID,
                'secret_key' => Str::random(6)
            ]);
        }

        return $record;
    }

    /**
     * 取消私有项目的分享
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return boolean 成功: true  失败: false
     */
    public static function closeShare($projectID, $userID)
    {
        $record = ProjectShare::where([
            ['project_id', $projectID],
            ['user_id', $userID]
        ])->first();

        if ($record) {
            return (bool)$record->delete();
        }

        return true;
    }

    /**
     * 修改分享的访问密码
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @return boolean|string 成功: 新的访问密码  失败: false
     */
    public static function changeSharePassword($projectID, $userID)
    {
        $secretKey = Str::random(6);

        if (!ProjectShare::where([['project_id', $projectID], ['user_id', $userID]])->update(['secret_key' => $secretKey])) {
            return false;
        }

        return $secretKey;
    }
}
