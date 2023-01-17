<?php

namespace App\Repositories\Iteration;

use App\Models\ProjectStar;
use Illuminate\Support\Facades\Auth;

class StarRepository
{
    /**
     * 用户是否收藏此项目
     *
     * @param int $userId 用户id
     * @param int $projectId 项目id
     * @return boolean
     */
    public static function hasStared(int $userId, int $projectId)
    {
        return (bool)ProjectStar::where([
            ['user_id', $userId],
            ['project_id', $projectId]
        ])->value('id');
    }

    /**
     * 收藏的项目id
     *
     * @param int $userId 用户id
     * @return array
     */
    public static function projectIds(int $userId)
    {
        return ProjectStar::where('user_id', $userId)->oldest('display_order')->pluck('project_id')->toArray();
    }

    /**
     * 收藏项目
     *
     * @param int $projectId 项目id
     * @return boolean
     */
    public static function create(int $projectId)
    {
        $exist = ProjectStar::where([
            ['user_id', Auth::id()],
            ['project_id', $projectId]
        ])->value('id');

        if ($exist) {
            return true;
        }

        $displayOrder = ProjectStar::where('user_id', Auth::id())->latest('display_order')->value('display_order');
        if (!$displayOrder) {
            $displayOrder = 1;
        } else {
            $displayOrder++;
        }

        ProjectStar::create([
            'user_id' => Auth::id(),
            'project_id' => $projectId,
            'display_order' => $displayOrder
        ]);

        return true;
    }

    /**
     * 修改顺序
     *
     * @param array $projectIds 项目id数组
     * @return void
     */
    public static function editOrder(array $projectIds)
    {
        $i = 1;
        foreach ($projectIds as $id) {
            ProjectStar::where([
                ['user_id', Auth::id()],
                ['project_id', $id]
            ])->update(['display_order' => $i]);

            $i++;
        }
    }

    /**
     * 取消收藏
     *
     * @param int $projectId 项目id
     * @return void
     */
    public static function del(int $projectId)
    {
        ProjectStar::where([
            ['user_id', Auth::id()],
            ['project_id', $projectId]
        ])->delete();
    }

    /**
     * 取消所有人的收藏
     *
     * @param int $projectId 项目id
     * @return void
     */
    public static function delAll(int $projectId)
    {
        ProjectStar::where('project_id', $projectId)->delete();
    }
}
