<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('iteration_apis', function (Blueprint $table) {
            $table->id();
            $table->bigInteger('iteration_id')->index()->comment('迭代id');
            $table->bigInteger('node_id')->comment('节点id');
            $table->tinyInteger('node_type')->default(0)->comment('节点类型:0目录,1文档');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('iteration_apis');
    }
};
