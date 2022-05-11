<?php

namespace App\Repositories\Project;

use App\Models\ProjectGroup;
use App\Models\GroupToProject;
use Illuminate\Database\Eloquent\Collection;

class GroupRepository
{
    /**
     * 获取分组列表
     * @param int $userID 用户id
     * @return Collection
     */
    public static function list($userID)
    {
        return ProjectGroup::where('user_id', $userID)->oldest('display_order')->get();
    }

    /**
     * 添加项目分组
     * @param int $userID 团队成员id
     * @param string $name 分组名称
     * @return ProjectGroup
     */
    public static function create($userID, $name)
    {
        $displayOrder = ProjectGroup::where('user_id', $userID)->latest('display_order')->value('display_order');
        if (!$displayOrder) {
            $displayOrder = 1;
        } else {
            $displayOrder++;
        }

        return ProjectGroup::create([
            'user_id' => $userID,
            'name' => $name,
            'display_order' => $displayOrder
        ]);
    }

    /**
     * 检查分组名称是否存在
     * 
     * @param int $userID 用户id
     * @param string $name 分组名称
     * @return boolean
     */
    public static function checkNameExist($userID, $name)
    {
        return (bool)ProjectGroup::where([
            ['user_id', $userID],
            ['name', $name]
        ])->value('id');
    }

    /**
     * 通过id获取分组信息
     * @param int $id 分组id
     * @return ProjectGroup
     */
    public static function find($id)
    {
        return ProjectGroup::find($id);
    }

    /**
     * 修改分组名称
     *
     * @param int $id 分组id
     * @param string $name 分组名称
     * @return boolean
     */
    public static function editName($id, $name)
    {
        return (bool)ProjectGroup::where('id', $id)->update(['name' => $name]);
    }

    /**
     * 修改分组顺序
     * @param array $ids 分组id数组
     * @return void
     */
    public static function editOrder($ids)
    {
        $i = 1;
        foreach ($ids as $id) {
            ProjectGroup::where('id', $id)->update(['display_order' => $i]);
            $i++;
        }
    }

    /**
     * 删除分组
     * @param int $id 分组id
     * @return boolean
     */
    public static function del($id)
    {
        return (bool)ProjectGroup::where('id', $id)->delete();
    }

    /**
     * 获取一个分组下的所有项目id
     *
     * @param int $id 分组id
     * @return array
     */
    public static function projectIds($id)
    {
        return GroupToProject::where('group_id', $id)->pluck('project_id')->toArray();
    }

    /**
     * 修改项目分组
     *
     * @param int $projectID 项目id
     * @param int $oldGroupID 原分组id
     * @param int $newGroupID 新分组id
     * @return boolean
     */
    public static function changeGroup($projectID, $oldGroupID, $newGroupID)
    {
        if ($oldGroupID > 0) {
            // 项目已在分组中
            $group = GroupToProject::where([
                ['group_id', $oldGroupID],
                ['project_id', $projectID]
            ])->first();

            if (!$group) {
                return false;
            }

            if ($newGroupID > 0) {
                // 项目从老分组移动到新分组
                $group->group_id = $newGroupID;
                return (bool)$group->save();
            } else {
                // 项目取消分组
                return (bool)$group->delete();
            }
        } else {
            // 项目当前没有分组
            if ($newGroupID > 0) {
                $group = GroupToProject::create([
                    'group_id' => $newGroupID,
                    'project_id' => $projectID
                ]);

                return (bool)$group;
            }

            return false;
        }
    }

    /**
     * 获取分组和项目对应关系
     *
     * @param int $id 分组id
     * @param int $userID 成员id
     * @param string $key 返回的数组key，只能是两种关系project => group或group => project
     * @return array
     */
    public static function relationship($id = 0, $userID = 0, $key = 'project')
    {
        if ($key != 'project' and $key != 'group') {
            return [];
        }

        if ($id) {
            if ($key == 'project') {
                return GroupToProject::where('group_id', $id)->pluck('group_id', 'project_id')->toArray();
            } else {
                return GroupToProject::where('group_id', $id)->pluck('project_id', 'group_id')->toArray();
            }
        }

        if (!$userID) {
            return [];
        }

        $groupIds = ProjectGroup::where('user_id', $userID)->pluck('id')->toArray();
        if (!$groupIds) {
            return [];
        }

        if ($key == 'project') {
            return GroupToProject::whereIn('group_id', $groupIds)->pluck('group_id', 'project_id')->toArray();
        } else {
            return GroupToProject::whereIn('group_id', $groupIds)->pluck('project_id', 'group_id')->toArray();
        }
    }
}