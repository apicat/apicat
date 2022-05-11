<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Repositories\ApiDoc\MockPathRepository;
use App\Repositories\Project\ApiDocRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\Project\TrashRepository;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class DocTrashController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    /**
     * 文档列表
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request)
    {
        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有查看回收站的权限。',
            ]);
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => TrashRepository::apiDocs(ProjectRepository::active()->id)
        ];
    }

    /**
     * 恢复已删除的api文档
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function restoreApiDoc(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
            'node_id' => ['nullable', 'integer', 'min:0']
        ]);

        if (!ProjectRepository::active()->hasAuthority()) {
            throw ValidationException::withMessages([
                'project_id' => '抱歉，您没有恢复该文档的权限。',
            ]);
        }

        if ($request->has('node_id')) {
            if ($request->input('node_id') > 0) {
                $node = ApiDocRepository::getNode($request->input('node_id'));
                if (!$node or ProjectRepository::active()->id != $node->project_id) {
                    throw ValidationException::withMessages([
                        'node_id' => '分类不存在',
                    ]);
                }
            }

            $res = TrashRepository::apiDocRestore($request->input('project_id'), $request->input('doc_id'), $request->input('node_id'));
        } else {
            if (!TrashRepository::checkNodeExist($request->input('doc_id'))) {
                return ['status' => -102, 'msg' => '原文档分类已被删除'];
            }
            $res = TrashRepository::apiDocRestore($request->input('project_id'), $request->input('doc_id'));
        }

        if (!$res) {
            throw ValidationException::withMessages([
                'result' => '文档恢复失败，请稍后重试。',
            ]);
        }

        MockPathRepository::restore(ProjectRepository::active()->id, $request->input('doc_id'));

        return ['status' => 0, 'msg' => '文档恢复成功'];
    }
}
