<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\Project\ApiDocRepository;
use App\Repositories\Project\ProjectRepository;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class DirectoryController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 所有分类列表
     * @param Request $request
     * @return array
     */
    public function list(Request $request)
    {
        return [
            'status' => 0,
            'msg' => '',
            'data' => ApiDocRepository::getDirTree($request->input('project_id'))
        ];
    }

    /**
     * 创建分类
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function store(Request $request)
    {
        $request->validate([
            'name' => ['required', 'string', 'max:50'],
            'parent_id' => ['nullable', 'integer', 'min:1']
        ]);

        $parentID = $request->input('parent_id') ? $request->input('parent_id') : 0;

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法创建目录',
            ]);
        }

        if ($node = ApiDocRepository::addDirToHead($request->input('project_id'), $request->input('name'), $parentID)) {
            return [
                'status' => 0,
                'msg' => '目录创建成功',
                'data' => [
                    'id' => $node->id,
                    'parent_id' => $node->parent_id,
                    'title' => $node->title,
                    'doc_id' => 0,
                    'sub_nodes' => []
                ]
            ];
        }

        throw ValidationException::withMessages([
            'result' => '目录创建失败，请稍后重试。',
        ]);
    }

    /**
     * 分类重命名
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function rename(Request $request)
    {
        $request->validate([
            'node_id' => ['required', 'integer', 'min:1'],
            'title' => ['required', 'string', 'max:50']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法重命名目录',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('node_id'));
        if (!$node or $node->project_id != $request->input('project_id')) {
            throw ValidationException::withMessages([
                'node_id' => '目录信息有误，无法重命名。',
            ]);
        }

        if (ApiDocRepository::rename($node, $request->input('title'))) {
            return [
                'status' => 0,
                'msg' => '目录重命名成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '目录重命名失败，请稍后重试。',
        ]);
    }

    /**
     * 删除分类
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function remove(Request $request)
    {
        $request->validate([
            'node_id' => ['required', 'integer', 'min:1']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '您无法删除目录',
            ]);
        }

        $node = ApiDocRepository::getNode($request->input('node_id'));
        if (!$node or $node->project_id != $request->input('project_id')) {
            throw ValidationException::withMessages([
                'node_id' => '目录信息有误，无法删除。',
            ]);
        }

        if (ApiDocRepository::removeDir($request->input('project_id'), $request->input('node_id'))) {
            // 该目录之后的树节点向前移动一位
            ApiDocRepository::moveForward($request->input('project_id'), $node->parent_id, $node->display_order);

            return [
                'status' => 0,
                'msg' => '目录删除成功'
            ];
        }

        throw ValidationException::withMessages([
            'result' => '目录删除失败，请稍后重试。',
        ]);
    }
}
