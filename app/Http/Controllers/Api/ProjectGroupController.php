<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Project\GroupRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Validation\ValidationException;

class ProjectGroupController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api']);
    }

    /**
     * 分组列表
     * @param Request $request
     * @return array
     */
    public function projectGroups(Request $request)
    {
        $groupArr = [];
        $groups = GroupRepository::list(Auth::id());

        if ($groups->count() > 0) {
            foreach ($groups as $group) {
                $groupArr[] = [
                    'id' => $group->id,
                    'name' => $group->name
                ];
            }
        }

        return ['status' => 0, 'msg' => '', 'data' => $groupArr];
    }

    /**
     * 创建项目
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function createGroup(Request $request)
    {
        $request->validate([
            'name' => 'required|string|max:255'
        ]);

        if (GroupRepository::checkNameExist(Auth::id(), $request->input('name'))) {
            throw ValidationException::withMessages([
                'name' => '该分组名称已经存在',
            ]);
        }

        $group = GroupRepository::create(Auth::id(), $request->input('name'));
        if (!$group) {
            throw ValidationException::withMessages([
                'result' => '添加分组失败，请稍后重试。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '添加分组成功',
            'data' => [
                'id' => $group->id,
                'name' => $group->name
            ]
        ];
    }

    /**
     * 重命名分组
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeName(Request $request)
    {
        $request->validate([
            'id' => 'required|integer|min:1',
            'name' => 'required|string|max:255'
        ]);

        $group = GroupRepository::find($request->input('id'));
        if (!$group or $group->user_id != Auth::id()) {
            throw ValidationException::withMessages([
                'id' => '分组不存在',
            ]);
        }

        if ($request->input('name') == $group->name) {
            return [
                'status' => 0,
                'msg' => '分组名称修改成功'
            ];
        }

        if (GroupRepository::checkNameExist(Auth::id(), $request->input('name'))) {
            throw ValidationException::withMessages([
                'name' => '该分组名称已经存在',
            ]);
        }

        if (GroupRepository::editName($group->id, $request->input('name'))) {
            return [
                'status' => 0,
                'msg' => '分组名称修改成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '分组名称修改失败，请稍后重试。',
        ]);
    }

    /**
     * 分组排序
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function changeOrder(Request $request)
    {
        $request->validate([
            'ids' => 'required|array|min:1'
        ]);

        $ids = $request->input('ids');
        foreach ($ids as $id) {
            if (!is_numeric($id)) {
                throw ValidationException::withMessages([
                    'ids' => '参数传递有误',
                ]);
            }
        }

        if (count(array_unique($ids)) != count($request->input('ids'))) {
            // id数组中有id重复
            throw ValidationException::withMessages([
                'ids' => '参数传递有误',
            ]);
        }

        GroupRepository::editOrder($ids);
        return ['status' => 0, 'msg' => '分组顺序修改成功'];
    }

    /**
     * 删除分组
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function remove(Request $request)
    {
        $request->validate([
            'id' => 'required|integer|min:1'
        ]);

        $group = GroupRepository::find($request->input('id'));
        if (!$group or $group->user_id != Auth::id()) {
            throw ValidationException::withMessages([
                'id' => '分组不存在',
            ]);
        }

        if (GroupRepository::del($group->id)) {
            return ['status' => 0, 'msg' => '分组删除成功'];
        }

        throw ValidationException::withMessages([
            'result' => '分组删除失败，请稍后重试。',
        ]);
    }
}
