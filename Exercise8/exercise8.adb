with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

procedure exercise8 is

    Count_Failed    : exception;    -- Exception to be raised when counting fails
    Gen             : Generator;    -- Random number generator

    protected type Transaction_Manager (N : Positive) is
        entry Finished;
        entry Wait_Until_Aborted; 
        function Commit return Boolean;
        procedure Signal_Abort;
    private
        Finished_Gate_Open  : Boolean := False;
        Aborted             : Boolean := False;
        Should_Commit       : Boolean := True;
    end Transaction_Manager;
    protected body Transaction_Manager is
      

        entry Wait_Until_Aborted when Aborted is
            begin
                if Wait_Until_Aborted'Count = N-3 then
                    Aborted := False;
                end if;
        end Wait_Until_Aborted; 


        entry Finished when Finished_Gate_Open or Finished'Count = N is
            begin
                if  Finished'Count = N-1 then
                    Finished_Gate_Open := True;           
                elsif Finished'Count = N-3 then 
                    Finished_Gate_Open := False;                          
                end if; 
                
            end Finished;




        procedure Signal_Abort is
        begin
            Aborted := True;
        end Signal_Abort;

        function Commit return Boolean is
        begin
            return Should_Commit;
        end Commit;
        
    end Transaction_Manager;



    
    function Unreliable_Slow_Add (x : Integer) return Integer is
    Error_Rate : Constant := 0.15;  -- (between 0 and 1)
    begin
-- PART 1
        if Random(Gen) < Error_Rate  then 
            delay Duration(0.5*Random(Gen));
            raise Count_Failed; 
        else 
            delay Duration(Random(Gen)*4.0); 
            return x + 10; 
        end if; 

 
    end Unreliable_Slow_Add;




    task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
    task body Transaction_Worker is
        Num         : Integer   := Initial;
        Prev        : Integer   := Num;
        Round_Num   : Integer   := 0;
    begin
        Put_Line ("Worker" & Integer'Image(Initial) & " started");

        loop
            Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
            Round_Num := Round_Num + 1;
-- PART 2
            select 
                Manager.Wait_Until_Aborted; 
                Num := Num + 5; 
            
            then abort 

                begin
                    Num := Unreliable_Slow_Add(Num); 
                exception 
                when Count_Failed => Manager.Signal_Abort; 
                end; 
                Manager.Finished; 
            end select; 

            
            Put_Line ("  Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num));
            
            Prev := Num;
            delay 0.5;

        end loop;
    end Transaction_Worker;

    Manager : aliased Transaction_Manager (3);

    Worker_1 : Transaction_Worker (0, Manager'Access);
    Worker_2 : Transaction_Worker (1, Manager'Access);
    Worker_3 : Transaction_Worker (2, Manager'Access);

begin
    Reset(Gen); -- Seed the random number generator
end exercise8;
