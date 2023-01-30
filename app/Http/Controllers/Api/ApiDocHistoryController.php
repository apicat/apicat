<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Modules\EditorJsonToHtml\Parser;
use App\Repositories\ApiDoc\ApiDocHistoryRepository;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\User\UserRepository;
use Illuminate\Validation\ValidationException;
use Illuminate\Http\Request;

class ApiDocHistoryController extends Controller
{
    public function __construct()
    {
        $this->middleware(['auth:api', 'in.this.project']);
    }

    public function histories(Request $request)
    {
        $request->validate([
            'doc_id' => ['required', 'integer', 'min:1'],
        ]);

        $records = ApiDocHistoryRepository::list($request->input('doc_id'));
        if ($records->isEmpty()) {
            return [
                'status' => 0,
                'msg' => '',
                'data' => []
            ];
        }

        $users = UserRepository::idNameArr(true);

        $result = [];
        foreach ($records as $record) {
            $month = $record->last_updated_at->format('Y-m');
            if (!isset($result[$month])) {
                $result[$month] = [
                    [
                        'id' => $record->id,
                        'title' => $record->last_updated_at->format('m月d日 H:i') . '(' . (isset($users[$record->last_user_id]) ? $users[$record->last_user_id] : '') . ')',
                        'type' => 1
                    ]
                ];
            } else {
                $result[$month][] = [
                    'id' => $record->id,
                    'title' => $record->last_updated_at->format('m月d日 H:i') . '(' . (isset($users[$record->last_user_id]) ? $users[$record->last_user_id] : '') . ')',
                    'type' => 1
                ];
            }
        }

        $result2 = [];
        foreach ($result as $k => $v) {
            $result2[] = [
                'id' => 0,
                'title' => str_replace('-', '年', $k) . '月',
                'type' => 0,
                'sub_nodes' => $v
            ];
        }

        return [
            'status' => 0,
            'msg' => '',
            'data' => $result2
        ];
    }

    public function detail(Request $request)
    {
        $request->validate([
            'id' => ['required', 'integer', 'min:1'],
        ]);

        if (!$record = ApiDocHistoryRepository::get($request->input('id'))) {
            throw ValidationException::withMessages([
                'id' => '您访问的历史记录不存在',
            ]);
        }

        if (!ApiDocRepository::inThisProject(ProjectRepository::active()->id, $record->doc_id)) {
            throw ValidationException::withMessages([
                'id' => '您访问的历史记录不存在',
            ]);
        }

        $content = $record->content ? Parser::parse($record->content, ProjectRepository::active()->id, $record->doc_id) : '';

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'id' => $record->id,
                'doc_id' => $record->doc_id,
                'title' => $record->title,
                'content' => $content,
                'created_time' => $record->last_updated_at->format('Y-m-d H:i'),
                'last_updated_by' => UserRepository::name($record->last_user_id, true)
            ]
        ];
    }

    public function diff(Request $request)
    {
        $request->validate([
            'id1' => ['required', 'integer', 'min:1'],
            'id2' => ['required', 'integer', 'min:0'],
        ]);

        if (!$record1 = ApiDocHistoryRepository::get($request->input('id1'))) {
            throw ValidationException::withMessages([
                'id1' => '您访问的历史记录不存在',
            ]);
        }
        if (!ApiDocRepository::inThisProject(ProjectRepository::active()->id, $record1->doc_id)) {
            throw ValidationException::withMessages([
                'id1' => '您访问的历史记录不存在',
            ]);
        }

        if ($request->input('id2')) {
            if (!$record2 = ApiDocHistoryRepository::get($request->input('id2'))) {
                throw ValidationException::withMessages([
                    'id2' => '您访问的历史记录不存在',
                ]);
            }
            if (!ApiDocRepository::inThisProject(ProjectRepository::active()->id, $record2->doc_id)) {
                throw ValidationException::withMessages([
                    'id2' => '您访问的历史记录不存在',
                ]);
            }
        } else {
            // id2值为0去查当前最新的文档内容
            if (!$record2 = ApiDocRepository::getNode($record1->doc_id)) {
                throw ValidationException::withMessages([
                    'id2' => '您访问的历史记录不存在',
                ]);
            }
            if (ProjectRepository::active()->id != $record2->project_id) {
                throw ValidationException::withMessages([
                    'id2' => '您访问的历史记录不存在',
                ]);
            }

            $record2->doc_id = $record2->id;
            $record2->title = $record2->title;
            $record2->last_updated_at = $record2->updated_at;
            $record2->last_user_id = $record2->updated_user_id;
        }

        $content1 = $record1->content ? Parser::parse($record1->content, ProjectRepository::active()->id, $record1->doc_id) : '';
        $content2 = $record2->content ? Parser::parse($record2->content, ProjectRepository::active()->id, $record2->doc_id) : '';

        return [
            'status' => 0,
            'msg' => '',
            'data' => [
                'doc1' => [
                    'id' => $record1->id,
                    'doc_id' => $record1->doc_id,
                    'title' => $record1->title,
                    'content' => $content1,
                    'created_time' => $record1->last_updated_at->format('Y-m-d H:i'),
                    'last_updated_by' => UserRepository::name($record1->last_user_id, true)
                ],
                'doc2' => [
                    'id' => $record2->id,
                    'doc_id' => $record2->doc_id,
                    'title' => $record2->title,
                    'content' => $content2,
                    'created_time' => $record2->last_updated_at->format('Y-m-d H:i'),
                    'last_updated_by' => UserRepository::name($record2->last_user_id, true)
                ]
            ]
        ];
    }
}
